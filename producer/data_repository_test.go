// le puse dentro de la estrucutra en mayuscula pero hay q ponerlo en minuscula pq asi va en la base de datos
package producer

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mock de la Base de Datos ---
// MockClient es una implementación simulada y completa de la interfaz mongoapi.Client.
// Debe implementar todos los métodos de la interfaz para ser un sustituto válido.
type MockClient struct {
	// Almacenamos los datos que queremos que devuelva nuestro mock.
	// La clave del mapa es el nombre de la colección y el valor son los datos.
	DbData map[string]interface{}
}

// RestfulAPIGetOne es la implementación simulada que usaremos.
func (mc *MockClient) RestfulAPIGetOne(collName string, filter bson.M) (map[string]interface{}, error) {
	// Devuelve los datos predefinidos para la colección solicitada.
	data, ok := mc.DbData[collName]
	if !ok {
		// Si no hay datos para esa colección, devuelve nil, simulando que no se encontró el documento.
		return nil, nil
	}
	return data.(map[string]interface{}), nil
}

// --- Métodos no utilizados de la interfaz ---
// Implementamos el resto de los métodos de la interfaz mongoapi.Client con funciones vacías,
// ya que no son necesarios para esta prueba específica, pero son requeridos para que nuestro
// MockClient sea un reemplazo válido del cliente real.
func (mc *MockClient) RestfulAPIGetMany(collName string, filter bson.M) ([]map[string]interface{}, error) {
	return nil, nil
}
func (mc *MockClient) RestfulAPIPutOne(collName string, filter bson.M, putData map[string]interface{}) (bool, error) {
	return false, nil
}
func (mc *MockClient) RestfulAPIPutOneNotUpdate(collName string, filter bson.M, putData map[string]interface{}) (bool, error) {
	return false, nil
}
func (mc *MockClient) RestfulAPIDeleteOne(collName string, filter bson.M) error {
	return nil
}
func (mc *MockClient) RestfulAPIDeleteMany(collName string, filter bson.M) error {
	return nil
}
func (mc *MockClient) RestfulAPIJSONPatch(collName string, filter bson.M, patchData []byte) error {
	return nil
}
func (mc *MockClient) RestfulAPIJSONPatchExtend(collName string, filter bson.M, patchData []byte, field string) error {
	return nil
}
func (mc *MockClient) RestfulAPIMergePatch(collName string, filter bson.M, patchData map[string]interface{}) error {
	return nil
}

func (db *MockClient) RestfulAPIPost(collName string, filter bson.M, postData map[string]interface{}) (bool, error) {
	return false, nil
}

func (db *MockClient) RestfulAPIPostMany(collName string, filter bson.M, postDataArray []interface{}) error {
	return nil
}

func (db *MockClient) RestfulAPIPutMany(collName string, filterArray []primitive.M, putDataArray []map[string]interface{}) error {
	return nil
}

func (db *MockClient) RestfulAPIPutOneTimeout(collName string, filter bson.M, putData map[string]interface{}, timeout int32, timeField string) bool {
	return false
}

// --- Función de Prueba ---
func TestQueryAuthSubsDataProcedure(t *testing.T) {
	// --- PREPARACIÓN ---
	ueId := "imsi-208930100007488"
	collName := "subscriptionData.authenticationData.authenticationSubscription"
	k4CollName := "encription.keydata.k4"
	expectedK4Value := "test_k4_value_decrypted"

	// --- CASO DE PRUEBA 1: Éxito ---
	// Prueba el flujo donde se encuentran los datos de autenticación y la clave K4.
	t.Run("Success case - data and K4 found", func(t *testing.T) {
		// Preparamos el mock de la base de datos con los datos para este caso.
		// Usamos la estructura de datos que proporcionaste.
		mockDb := &MockClient{
			DbData: map[string]interface{}{
				// Datos para la colección de suscripciones de autenticación.
				collName: map[string]interface{}{
					"ueId":                 ueId,
					"authenticationMethod": "5G_AKA",
					"permanentKey": map[string]interface{}{
						"encryptionAlgorithm": float64(0), // los números en JSON se decodifican como float64
						"encryptionKey":       "0",        // Este valor será reemplazado
						"permanentKeyValue":   "5122250214c33e723a5dd523fc145fc0",
					},
					"sequenceNumber": "16f3b3f70fc2",
					"k4_sno":         float64(1), // Usamos float64 como en tu ejemplo
				},
				// Datos para la colección de claves K4.
				k4CollName: map[string]interface{}{
					"k4_sno": float64(1),
					"k4":     expectedK4Value,
				},
			},
		}
		// Reemplazamos el cliente de base de datos real por nuestro mock.
		// Esto es seguro hacerlo en una prueba.
		AuthDBClient = mockDb

		// --- EJECUCIÓN ---
		authData, problemDetails := QueryAuthSubsDataProcedure(collName, ueId)

		// --- VERIFICACIÓN ---
		assert.Nil(t, problemDetails)
		assert.NotNil(t, authData)
		assert.Equal(t, ueId, authData["ueId"])

		// Verificamos que la clave de encriptación ha sido reemplazada correctamente.
		permanentKey, ok := authData["permanentKey"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, expectedK4Value, permanentKey["encryptionKey"])
	})

	// --- CASO DE PRUEBA 2: Usuario no encontrado ---
	t.Run("User not found", func(t *testing.T) {
		// Usamos un mock de base de datos vacío para simular que el usuario no existe.
		mockDb := &MockClient{DbData: map[string]interface{}{}}
		AuthDBClient = mockDb

		// --- EJECUCIÓN ---
		authData, problemDetails := QueryAuthSubsDataProcedure(collName, "non-existent-ue")

		// --- VERIFICACIÓN ---
		assert.Nil(t, authData)
		assert.NotNil(t, problemDetails)
		assert.Equal(t, http.StatusNotFound, int(problemDetails.Status))
		assert.Equal(t, "USER_NOT_FOUND", problemDetails.Cause)
	})

	// --- CASO DE PRUEBA 3: Clave K4 no encontrada ---
	t.Run("K4 key not found", func(t *testing.T) {
		// Preparamos el mock solo con los datos de autenticación, pero sin la clave K4.
		mockDb := &MockClient{
			DbData: map[string]interface{}{
				collName: map[string]interface{}{
					"ueId":                 ueId,
					"authenticationMethod": "5G_AKA",
					"permanentKey": map[string]interface{}{
						"encryptionAlgorithm": float64(0),
						"encryptionKey":       "0", // Este valor NO será reemplazado
						"permanentKeyValue":   "5122250214c33e723a5dd523fc145fc0",
					},
					"k4_sno": float64(2), // Un SNO que no encontrará en la colección k4
				},
			},
		}
		AuthDBClient = mockDb

		// --- EJECUCIÓN ---
		authData, problemDetails := QueryAuthSubsDataProcedure(collName, ueId)

		// --- VERIFICACIÓN ---
		assert.Nil(t, problemDetails) // La función no debe devolver un error fatal.
		assert.NotNil(t, authData)

		// Verificamos que la clave de encriptación NO haya sido reemplazada.
		permanentKey, ok := authData["permanentKey"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "0", permanentKey["encryptionKey"]) // Se mantiene el valor original.
	})
}
