package svcregistry

import (
	"context"
	"testing"
	"time"
	"utils/logging"
	"utils/utils/testutils"
)

func TestConsul(t *testing.T) {

	l := logging.New()

	ctx := context.Background()

	consulAddr := "localhost"
	consulPort := 8500

	// pra rodar o consul localmente:
	// sudo docker compose -f 'consul.yml' up -d --build

	// isso deve ser executado antes dos testes abaixo.
	// pra matar o container:
	// sudo docker stop consul && sudo docker rm consul

	t.Run("reset", func(t *testing.T) {

		nowDT := time.Now()

		getter, err := NewConsulGetter(
			"localhost", 8500, 3*time.Second,
		)
		testutils.AssertBool(t, err == nil, true)

		///////
		// cache vazio

		getter.Reset(l, "baba")

		want := map[string]cacheControl{}
		testutils.AssertStruct(t, getter.cachedAddresses, want)

		///////
		// cache com um elemento

		getter.cachedAddresses["baba"] = cacheControl{
			creationDT: nowDT,
			cursor:     0,
			addresses:  []string{"lolo:1111"},
			// addresses: []ServiceAddress{
			// 	{
			// 		Address: "lolo",
			// 		Port:    1111,
			// 	},
			// },
		}

		getter.Reset(l, "baba")

		want = map[string]cacheControl{}
		testutils.AssertStruct(t, getter.cachedAddresses, want)

		///////
		// cache com 2 elementos, cursor no 0

		getter.cachedAddresses["baba"] = cacheControl{
			creationDT: nowDT,
			cursor:     0,
			addresses: []string{
				"lolo:1111",
				"lala:1122",
			},
			// addresses: []ServiceAddress{
			// 	{
			// 		Address: "lolo",
			// 		Port:    1111,
			// 	},
			// 	{
			// 		Address: "lala",
			// 		Port:    1122,
			// 	},
			// },
		}

		getter.Reset(l, "baba")

		// só pode ficar o elemento 0
		want = map[string]cacheControl{
			"baba": {
				creationDT: nowDT,
				cursor:     0,
				addresses: []string{
					"lolo:1111",
				},
				// addresses: []ServiceAddress{
				// 	{
				// 		Address: "lolo",
				// 		Port:    1111,
				// 	},
				// },
			},
		}
		testutils.AssertStruct(t, getter.cachedAddresses, want)

		///////
		// cache com 2 elementos, cursor no 1

		getter.cachedAddresses["baba"] = cacheControl{
			creationDT: nowDT,
			cursor:     1,
			addresses: []string{
				"lolo:1111",
				"lala:1122",
			},
			// addresses: []ServiceAddress{
			// 	{
			// 		Address: "lolo",
			// 		Port:    1111,
			// 	},
			// 	{
			// 		Address: "lala",
			// 		Port:    1122,
			// 	},
			// },
		}

		getter.Reset(l, "baba")

		// só pode ficar o elemento 1
		want = map[string]cacheControl{
			"baba": {
				creationDT: nowDT,
				cursor:     0,
				addresses: []string{
					"lala:1122",
				},
				// addresses: []ServiceAddress{
				// 	{
				// 		Address: "lala",
				// 		Port:    1122,
				// 	},
				// },
			},
		}
		testutils.AssertStruct(t, getter.cachedAddresses, want)

	})

	t.Run("register, heartbeat, get address, reset, get address",
		func(t *testing.T) {

			// gerenciador que permite registro e heartbeat
			h, err := NewConsulHandler(
				consulAddr, uint(consulPort),
			)
			testutils.AssertBool(t, err == nil, true)

			// registrando uma instancia do serviço
			err = h.RegisterInstance(
				l,
				ctx,
				"baba",
				"123",
				"localhost",
				1122,
				nil,
				60*time.Second,
			)
			testutils.AssertBool(t, err == nil, true)

			defer h.DeregisterInstance(l, ctx, "123")

			err = h.SendHeartbeat(
				l, ctx, "baba", "123",
			)
			testutils.AssertBool(t, err == nil, true)

			// registrando outra instancia do mesmo servico
			err = h.RegisterInstance(
				l,
				ctx,
				"baba",
				"456",
				"localhost",
				1123,
				nil,
				60*time.Second,
			)
			testutils.AssertBool(t, err == nil, true)

			defer h.DeregisterInstance(l, ctx, "456")

			err = h.SendHeartbeat(
				l, ctx, "baba", "456",
			)
			testutils.AssertBool(t, err == nil, true)

			/////////////
			// agora, criando apenas o getter

			getter, err := NewConsulGetter(
				"localhost", 8500, 3*time.Second,
			)
			testutils.AssertBool(t, err == nil, true)

			// vai consultar o consul, ja q nao tem no cache
			l.Info("vai pegar do consul")
			got, err := getter.GetServiceAddress(
				l, ctx, "baba",
			)
			testutils.AssertBool(t, err == nil, true)

			testutils.AssertBool(
				t,
				got == "localhost:1122" ||
					got == "localhost:1123",
				true,
			)

			// agora, vai pegar do cache
			l.Info("vai pegar do cache")
			got, err = getter.GetServiceAddress(
				l, ctx, "baba",
			)
			testutils.AssertBool(t, err == nil, true)

			testutils.AssertBool(
				t,
				got == "localhost:1122" ||
					got == "localhost:1123",
				true,
			)

			// colocando um sleep pra o cache expirar
			time.Sleep(5 * time.Second)

			// vai pegar do consul novamente,
			// ja q o cache expirou
			l.Info("vai pegar do consul")
			got, err = getter.GetServiceAddress(
				l, ctx, "baba",
			)
			testutils.AssertBool(t, err == nil, true)

			testutils.AssertBool(
				t,
				got == "localhost:1122" ||
					got == "localhost:1123",
				true,
			)

			///////////
			// forçando invalidaçao do cache

			getter.Reset(l, "baba")

			// vai pegar do consul novamente,
			// ja q o cache foi limpado
			l.Info("vai pegar do consul")
			got, err = getter.GetServiceAddress(
				l, ctx, "baba",
			)
			testutils.AssertBool(t, err == nil, true)

			testutils.AssertBool(
				t,
				got == "localhost:1122" ||
					got == "localhost:1123",
				true,
			)
		})

	t.Run("register, heartbeat, deregister", func(t *testing.T) {

		h, err := NewConsulHandler(
			consulAddr, uint(consulPort),
		)
		testutils.AssertBool(t, err == nil, true)

		err = h.RegisterInstance(
			l,
			ctx,
			"baba",
			"123",
			"localhost",
			1122,
			nil,
			60*time.Second,
		)
		testutils.AssertBool(t, err == nil, true)

		defer h.DeregisterInstance(
			l, ctx, "123",
		)

		err = h.SendHeartbeat(
			l, ctx, "baba", "123",
		)
		testutils.AssertBool(t, err == nil, true)

		err = h.RegisterInstance(
			l,
			ctx,
			"baba",
			"456",
			"localhost",
			1123,
			nil,
			60*time.Second,
		)
		testutils.AssertBool(t, err == nil, true)

		defer h.DeregisterInstance(
			l, ctx, "456",
		)

		err = h.SendHeartbeat(
			l, ctx, "baba", "456",
		)
		testutils.AssertBool(t, err == nil, true)

	})

	t.Run("next address", func(t *testing.T) {

		h, err := NewConsulGetter(
			"localhost", 8500, 3*time.Second,
		)
		testutils.AssertBool(t, err == nil, true)

		////////////
		// empty addresses

		cControl := cacheControl{
			creationDT: time.Now(),
			cursor:     0,
			addresses:  []string{},
		}
		h.cachedAddresses["baba"] = cControl

		got := h.nextAddress("baba")
		want := ""

		testutils.AssertString(t, got, want)

		// pegando next de novo, só pra garantir q nao
		// haverá panic
		got = h.nextAddress("baba")

		testutils.AssertStruct(t, got, want)

		////////////
		// um address

		cControl = cacheControl{
			creationDT: time.Now(),
			cursor:     0,
			addresses: []string{
				"huhuhu:1111",
			},
			// addresses: []ServiceAddress{
			// 	{
			// 		Address: "huhuhu",
			// 		Port:    1111,
			// 	},
			// },
		}
		h.cachedAddresses["bebe"] = cControl

		got = h.nextAddress("bebe")
		want = "huhuhu:1111"
		// want = ServiceAddress{
		// 	Address: "huhuhu",
		// 	Port:    1111,
		// }

		testutils.AssertStruct(t, got, want)

		// pegando next de novo, só pra garantir q nao
		// haverá panic
		got = h.nextAddress("bebe")

		testutils.AssertStruct(t, got, want)

		////////////
		// dois addresses

		cControl = cacheControl{
			creationDT: time.Now(),
			cursor:     0,
			addresses: []string{
				"huhuhu:1111",
				"papapa:2222",
			},
			// addresses: []ServiceAddress{
			// 	{
			// 		Address: "huhuhu",
			// 		Port:    1111,
			// 	},
			// 	{
			// 		Address: "papapa",
			// 		Port:    2222,
			// 	},
			// },
		}
		h.cachedAddresses["bibi"] = cControl

		got = h.nextAddress("bibi")
		want = cControl.addresses[0]

		testutils.AssertStruct(t, got, want)

		// pegando next de novo: vai ser o 2º elemento
		got = h.nextAddress("bibi")
		want = cControl.addresses[1]

		testutils.AssertStruct(t, got, want)

		// pegando next de novo: vai ser o 1º elemento
		got = h.nextAddress("bibi")
		want = cControl.addresses[0]

		testutils.AssertStruct(t, got, want)

		// pegando next de novo: vai ser o 2º elemento
		got = h.nextAddress("bibi")
		want = cControl.addresses[1]

		testutils.AssertStruct(t, got, want)
	})

}
