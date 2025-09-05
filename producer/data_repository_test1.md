package producer

import (
	"net/http"
	"testing"

	"github.com/omec-project/udr/context"
	"github.com/omec-project/udr/factory"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// MockClient es una implementación simulada de la interfaz de la base de datos.
// Nos permite controlar lo que devuelven las llamadas a la base de datos durante la prueba.
type MockClient struct {
	// Almacenamos los datos que queremos que devuelva nuestro mock.
	// La clave del mapa es el nombre de la colección y el valor son los datos.
	DbData map[string]interface{}
}

// RestfulAPIGetOne es la implementación simulada de la función que obtiene un documento.
func (mc *MockClient) RestfulAPIGetOne(collName string, filter bson.M) (map[string]interface{}, error) {
	// Devuelve los datos predefinidos para la colección solicitada.
	// Si no hay datos para esa colección, devuelve nil (simulando que no se encontró el documento).
	data, ok := mc.DbData[collName]
	if !ok {
		return nil, nil
	}
	return data.(map[string]interface{}), nil
}

// Implementamos las otras funciones de la interfaz mongoapi.Client, pero no las necesitamos
// para esta prueba específica, por lo que pueden estar vacías.
func (mc *MockClient) RestfulAPIGetMany(collName string, filter bson.M) ([]map[string]interface{}, error) {
	return nil, nil
}
func (mc *MockClient) RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) (bool, error) {
	return false, nil
}

// ... (y así sucesivamente para los otros métodos de la interfaz)

// TestQueryAuthSubsDataProcedure es la función principal de nuestra prueba.
func TestQueryAuthSubsDataProcedure(t *testing.T) {
	// --- PREPARACIÓN ---

	// Inicializamos la configuración de UDR para la prueba.
	factory.InitConfigFactory("./factory/udr_config.yaml")
	context.Init()

	// Definimos los datos de prueba que simularemos que están en la base de datos.
	ueId := "imsi-208930000000001"
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	k4CollName := "encription.keydata.k4"
	expectedK4Value := "test_k4_value"

	// --- CASO DE PRUEBA 1: Éxito ---
	// Este caso prueba el flujo feliz donde se encuentran tanto los datos de autenticación como la clave K4.
	t.Run("Success", func(t *testing.T) {
		// Preparamos el mock de la base de datos con los datos para este caso.
		mockDb := &MockClient{
			DbData: map[string]interface{}{
				// Datos para la colección de suscripciones de autenticación.
				collName: map[string]interface{}{
					"ueId": ueId,
					"permanentKey": map[string]interface{}{
						"encryptionKey": "old_key", // Un valor que esperamos que sea reemplazado.
					},
					"k4_sno": byte(1), // El SNO para buscar la clave K4.
				},
				// Datos para la colección de claves K4.
				k4CollName: map[string]interface{}{
					"k4_sno": byte(1),
					"k4":     expectedK4Value,
				},
			},
		}
		// Reemplazamos el cliente de base de datos real por nuestro mock.
		AuthDBClient = mockDb

		// --- EJECUCIÓN ---
		// Llamamos a la función que queremos probar.
		authData, problemDetails := QueryAuthSubsDataProcedure(collName, ueId)

		// --- VERIFICACIÓN ---
		// Verificamos que no haya ocurrido ningún problema.
		assert.Nil(t, problemDetails)
		// Verificamos que los datos de autenticación no sean nulos.
		assert.NotNil(t, authData)
		// Verificamos que el ueId en la respuesta sea el correcto.
		assert.Equal(t, ueId, authData["ueId"])

		// Verificamos que la clave de encriptación ha sido reemplazada correctamente por el valor de K4.
		permanentKey, ok := authData["permanentKey"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, expectedK4Value, permanentKey["encryptionKey"])
	})

	// --- CASO DE PRUEBA 2: Usuario no encontrado ---
	// Este caso prueba qué sucede si el UE no existe en la base de datos.
	t.Run("UserNotFound", func(t *testing.T) {
		// Preparamos un mock de base de datos vacío.
		mockDb := &MockClient{
			DbData: map[string]interface{}{},
		}
		AuthDBClient = mockDb

		// --- EJECUCIÓN ---
		authData, problemDetails := QueryAuthSubsDataProcedure(collName, "non-existent-ue")

		// --- VERIFICACIÓN ---
		// Verificamos que no se devuelvan datos de autenticación.
		assert.Nil(t, authData)
		// Verificamos que se devuelva un error de "USER_NOT_FOUND".
		assert.NotNil(t, problemDetails)
		assert.Equal(t, http.StatusNotFound, int(problemDetails.Status))
		assert.Equal(t, "USER_NOT_FOUND", problemDetails.Cause)
	})

	// --- CASO DE PRUEBA 3: Clave K4 no encontrada ---
	// Este caso prueba qué sucede si los datos de autenticación existen pero la clave K4 no.
	t.Run("K4KeyNotFound", func(t *testing.T) {
		// Preparamos el mock solo con los datos de autenticación, pero no los de K4.
		mockDb := &MockClient{
			DbData: map[string]interface{}{
				collName: map[string]interface{}{
					"ueId": ueId,
					"permanentKey": map[string]interface{}{
						"encryptionKey": "old_key",
					},
					"k4_sno": byte(2), // Un SNO que no existirá en la colección k4.
				},
			},
		}
		AuthDBClient = mockDb

		// --- EJECUCIÓN ---
		authData, problemDetails := QueryAuthSubsDataProcedure(collName, ueId)

		// --- VERIFICACIÓN ---
		// Verificamos que no haya ningún error fatal.
		assert.Nil(t, problemDetails)
		// Verificamos que se devuelvan los datos de autenticación.
		assert.NotNil(t, authData)

		// Verificamos que la clave de encriptación NO haya sido reemplazada,
		// ya que la clave K4 no se encontró.
		permanentKey, ok := authData["permanentKey"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "old_key", permanentKey["encryptionKey"])
	})
}
