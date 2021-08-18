package clusterapi

import (
	"fmt"
	"net/http"

	"github.com/semi-technologies/weaviate/adapters/handlers/rest/state"
)

func Serve(appState *state.State) {
	port := appState.ServerConfig.Config.Cluster.DataBindPort
	if port <= 0 {
		port = 7946
	}

	appState.Logger.WithField("port", port).
		WithField("action", "cluster_api_startup").
		Debugf("serving cluster api on port %d", port)

	schema := newSchema(appState.SchemaManager.TxManager())

	mux := http.NewServeMux()
	mux.Handle("/schema/transactions/",
		http.StripPrefix("/schema/transactions/", schema.transactions()))
	mux.Handle("/", schema.index())
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}