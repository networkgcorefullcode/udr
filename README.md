<!--
SPDX-FileCopyrightText: 2025 Canonical Ltd
SPDX-FileCopyrightText: 2021 Open Networking Foundation <info@opennetworking.org>
Copyright 2019 free5GC.org

SPDX-License-Identifier: Apache-2.0
-->
[![Go Report Card](https://goreportcard.com/badge/github.com/omec-project/udr)](https://goreportcard.com/report/github.com/omec-project/udr)

# UDR

Implements 3gpp 29.504 specification. Provides service to PCF, UDM. UDR supports
SBI interface and any other network function can use the service.

## UDR block diagram
![UDR Block Diagram](/docs/images/README-UDR.png)

## Repository Structure

Below is a high-level view of the repository and its main components:
```
.
├── consumer                    # Contains logic for inter-NF communication. Implements discovery and management procedures to interact with other network functions.
│   ├── nf_management.go
│   └── nf_management_test.go
├── context                     # Manages global service context, runtime data, and in-memory state shared across the UDR components.
│   └── context.go
├── datarepository              # Implements 3GPP-defined APIs for accessing and managing user and subscription data. Includes endpoints for authentication, registration, session management, and exposure data. This is the core of the UDR service.
│   ├── api_access_and_mobility_data.go
│   ├── api_access_and_mobility_subscription_data_document.go
│   ├── api_amf3_gpp_access_registration_document.go
│   ├── api_amf_non3_gpp_access_registration_document.go
│   ├── api_amf_subscription_info_document.go
│   ├── api_authentication_data_document.go
│   ├── api_authentication_so_r_document.go
│   ├── api_authentication_status_document.go
│   ├── api_default.go
│   ├── api_event_amf_subscription_info_document.go
│   ├── api_event_exposure_data_document.go
│   ├── api_event_exposure_group_subscription_document.go
│   ├── api_event_exposure_group_subscriptions_collection.go
│   ├── api_event_exposure_subscription_document.go
│   ├── api_event_exposure_subscriptions_collection.go
│   ├── api_operator_specific_data_container_document.go
│   ├── api_parameter_provision_document.go
│   ├── api_pdu_session_management_data.go
│   ├── api_provisioned_data_document.go
│   ├── api_provisioned_parameter_data_document.go
│   ├── api_query_amf_subscription_info_document.go
│   ├── api_query_identity_data_by_supi_or_gpsi_document.go
│   ├── api_query_odb_data_by_supi_or_gpsi_document.go
│   ├── api_retrieval_of_shared_data.go
│   ├── api_sdm_subscription_document.go
│   ├── api_sdm_subscriptions_collection.go
│   ├── api_session_management_subscription_data.go
│   ├── api_smf_registration_document.go
│   ├── api_smf_registrations_collection.go
│   ├── api_smf_selection_subscription_data_document.go
│   ├── api_smsf3_gpp_registration_document.go
│   ├── api_smsf_non3_gpp_registration_document.go
│   ├── api_sms_management_subscription_data_document.go
│   ├── api_sms_subscription_data_document.go
│   ├── api_subs_to_nofify_collection.go
│   ├── api_subs_to_notify_document.go
│   ├── api_trace_data_document.go
│   └── routers.go
├── dev-container.ps1
├── dev-container.sh
├── DEV_README.md
├── Dockerfile
├── Dockerfile_dev
├── Dockerfile.fast
├── docs                        # Contains documentation assets, including UDR architecture and operational diagrams.
│   └── images
│       ├── README-UDR.png
│       └── README-UDR.png.license
├── factory                     # Provides configuration management utilities. Loads and parses YAML configuration files (e.g., udr_config.yaml) and initializes system parameters. 
│   ├── config.go
│   ├── factory.go
│   ├── udr_config_test.go
│   ├── udr_config_with_custom_webui_url.yaml
│   └── udr_config.yaml
├── go.mod
├── go.mod.license
├── go.sum
├── go.sum.license
├── LICENSES
│   └── Apache-2.0.txt
├── logger                      # Centralized logging system configuration for the UDR service. Defines logging levels and output format.
│   └── logger.go
├── Makefile
├── metrics                     # Exposes telemetry and performance metrics for monitoring the UDR’s operational status.
│   └── telemetry.go
├── nfregistration              # Handles registration and deregistration of the UDR with the Network Repository Function (NRF).
│   ├── nf_registration.go
│   └── nf_registration_test.go
├── NOTICE.txt
├── polling                     # Implements mechanisms to periodically refresh NF configuration and maintain synchronization with network state.
│   ├── nf_configuration.go
│   └── nf_configuration_test.go
├── producer                    # Contains the core UDR business logic, including data handling, repository adapters, and internal callbacks. Defines how data is stored, retrieved, and served to other NFs.
│   ├── callback
│   │   └── callback.go
│   ├── callback.go
│   ├── data_repository.go
│   ├── data_repository_test1.md
│   ├── data_repository_test2.md
│   ├── data_repository_test.go
│   └── db_adapter.go
├── README.md
├── service                     # Contains the service initialization code that bootstraps all UDR modules during startup.
│   └── init.go
├── Taskfile.yml
├── test-mirror.txt
├── udr.go
├── util                        # Provides helper functions for data conversion, context initialization, and general utility operations.
│   ├── convert.go
│   ├── init_context.go
│   └── util.go
├── VERSION
└── VERSION.license

16 directories, 84 files
```

## Configuration and Deployment

**Docker**

To build the container image:
```
task mod-start
task build
task docker-build-fast
```

**Kubernetes**

The standard deployment uses Helm charts from the Aether project. The version of the Chart can be found in the OnRamp repository in the `vars/main.yml` file.


## Quick Navigation

| Goal                                       | Location                                                                                                                                        | Description                                                                                                      |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| **Start or build the UDR service**         | [`udr.go`](./udr.go) / [`service/init.go`](./service/init.go)                                                                                   | Entry point and service initialization logic.                                                                    |
| **Configure UDR**                          | [`factory/udr_config.yaml`](./factory/udr_config.yaml)                                                                                          | Main configuration file defining UDR parameters, service addresses, and NRF registration info.                   |
| **Explore 3GPP UDR APIs**                  | [`datarepository/`](./datarepository/)                                                                                                          | Main API implementation for data access and management (authentication data, subscriptions, session data, etc.). |
| **Understand NF registration flow**        | [`nfregistration/`](./nfregistration/)                                                                                                          | Code for registering the UDR instance with the NRF.                                                              |
| **Modify data storage or retrieval logic** | [`producer/data_repository.go`](./producer/data_repository.go)                                                                                  | Core UDR data handling layer and business logic implementation.                                                  |
| **Enable logging or adjust log output**    | [`logger/logger.go`](./logger/logger.go)                                                                                                        | Configuration for UDR logging infrastructure.                                                                    |
| **Monitor metrics and telemetry**          | [`metrics/telemetry.go`](./metrics/telemetry.go)                                                                                                | Exposes metrics for Prometheus or other monitoring tools.                                                        |
| **Check context management**               | [`context/context.go`](./context/context.go)                                                                                                    | Manages runtime and service-level state for the UDR.                                                             |
| **Work with NF discovery/management**      | [`consumer/`](./consumer/)                                                                                                                      | Implements inter-NF management interfaces.                                                                       |
| **Run or inspect tests**                   | [`producer/data_repository_test.go`](./producer/data_repository_test.go) / [`consumer/nf_management_test.go`](./consumer/nf_management_test.go) | Unit tests for repository logic and NF management.                                                               |
| **View UDR architecture diagrams**         | [`docs/images/README-UDR.png`](./docs/images/README-UDR.png)                                                                                    | Visual overview of the UDR service and its interactions.                                                         |
| **Utility and helper functions**           | [`util/`](./util/)                                                                                                                              | Conversion and initialization utilities used throughout the codebase.                                            |



## Dynamic Network configuration (via webconsole)

UDR polls the webconsole every 5 seconds to fetch the latest PLMN configuration.

### Setting Up Polling

Include the `webuiUri` of the webconsole in the configuration file
```
configuration:
  ...
  webuiUri: https://webui:5001 # or http://webui:5001
  ...
```
The scheme (http:// or https://) must be explicitly specified. If no parameter is specified,
UDR will use `http://webui:5001` by default.

### HTTPS Support

If the webconsole is served over HTTPS and uses a custom or self-signed certificate,
you must install the root CA certificate into the trust store of the UDR environment.

Check the official guide for installing root CA certificates on Ubuntu:
[Install a Root CA Certificate in the Trust Store](https://documentation.ubuntu.com/server/how-to/security/install-a-root-ca-certificate-in-the-trust-store/index.html)

## Upcoming changes
- Subscription management callbacks to network functions.

Compliance of the 5G Network functions can be found at [5G Compliance](https://docs.sd-core.opennetworking.org/main/overview/3gpp-compliance-5g.html)

## Reach out to us through

1. #sdcore-dev channel in [ONF Community Slack](https://aether5g-project.slack.com)
2. Raise Github [issues](https://github.com/omec-project/udr/issues/new)
