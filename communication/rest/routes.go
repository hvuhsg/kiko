package rest

import (
	"net/http"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/google/uuid"
	"github.com/hvuhsg/kiko/execution"
)

func writeError(errorMessage string, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := simplejson.New()
	response.Set("error", errorMessage)

	payload, err := response.MarshalJSON()
	if err != nil {
		panic(err)
	}

	w.Write(payload)
}

func handleNode(w http.ResponseWriter, r *http.Request, engine *execution.IEngine) {
	switch r.Method {
	case "POST":
		// create node
		nodeUuid := (*engine).AddNode()

		// create json response
		json := simplejson.New()
		json.Set("uuid", nodeUuid.String())
		response, _ := json.MarshalJSON()

		// send response
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Location", r.URL.String()+"?node="+nodeUuid.String())
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	case "GET":
		// extraxt node uuid argument
		nodeUuidStr := r.URL.Query().Get("node")
		if nodeUuidStr == "" {
			writeError("Node UUID is required", w)
			return
		}

		// validate argument
		nodeUuid, err := uuid.Parse(nodeUuidStr)
		if err != nil {
			writeError("node argument is not a valid UUID", w)
			return
		}

		connections, err := (*engine).GetNodeConnections(nodeUuid)
		if err != nil {
			writeError(err.Error(), w)
			return
		}

		// create json response
		json := simplejson.New()
		json.Set("uuid", nodeUuid.String())
		json.Set("connections", connections)
		response, _ := json.MarshalJSON()

		// send response
		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	case "DELETE":
		// load node uuid argument
		nodeUuidStr := r.URL.Query().Get("node")
		if nodeUuidStr == "" {
			writeError("Node UUID is required", w)
			return
		}

		// validate argument
		nodeUuid, err := uuid.Parse(nodeUuidStr)
		if err != nil {
			writeError("node argument is not a valid UUID", w)
			return
		}

		// delete node
		err = (*engine).RemoveNode(nodeUuid)
		if err != nil {
			writeError(err.Error(), w)
			return
		}

		// send response
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
	default:
		writeError("Method not allowed", w)
	}
}

func handleConnection(w http.ResponseWriter, r *http.Request, engine *execution.IEngine) {
	switch r.Method {
	case "POST":
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")
		weight := r.URL.Query().Get("weight")

		// Validate arguments
		if from == "" {
			writeError("Connection requires 'from' UUID argument", w)
			return
		}
		if to == "" {
			writeError("Connection requires 'to' UUID argument", w)
			return
		}
		if weight == "" {
			writeError("Connection requires 'weight' int argument", w)
			return
		}

		fromUuid, err := uuid.Parse(from)
		if err != nil {
			writeError("'from' argument is invalid UUID", w)
			return
		}

		toUuid, err := uuid.Parse(to)
		if err != nil {
			writeError("'to' argument is invalid UUID", w)
			return
		}

		weightInt, err := strconv.ParseUint(weight, 10, 0)
		if err != nil {
			writeError("'weight' argument is invalid int", w)
			return
		}

		err = (*engine).AddConnection(fromUuid, toUuid, uint(weightInt))
		if err != nil {
			writeError(err.Error(), w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	case "DELETE":
		from := r.URL.Query().Get("from")
		to := r.URL.Query().Get("to")

		// Validate arguments
		if from == "" {
			writeError("Connection requires 'from' UUID argument", w)
			return
		}
		if to == "" {
			writeError("Connection requires 'to' UUID argument", w)
			return
		}

		fromUuid, err := uuid.Parse(from)
		if err != nil {
			writeError("'from' argument is invalid UUID", w)
			return
		}

		toUuid, err := uuid.Parse(to)
		if err != nil {
			writeError("'to' argument is invalid UUID", w)
			return
		}

		err = (*engine).RemoveConnection(fromUuid, toUuid)
		if err != nil {
			writeError(err.Error(), w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	default:
		writeError("Method not allowed", w)
	}
}

func handleKnnQuery(w http.ResponseWriter, r *http.Request, engine *execution.IEngine) {
	if r.Method != "GET" {
		writeError("Method not allowed", w)
		return
	}

	nodeUuidStr := r.URL.Query().Get("node")
	kStr := r.URL.Query().Get("k")

	if nodeUuidStr == "" {
		writeError("Node UUID is required", w)
		return
	}

	if kStr == "" {
		writeError("K int is required", w)
		return
	}

	nodeUuid, err := uuid.Parse(nodeUuidStr)
	if err != nil {
		writeError("node argument is not a valid UUID", w)
		return
	}

	kInt, err := strconv.ParseUint(kStr, 10, 0)
	if err != nil {
		writeError("k argument is not a vaild uint", w)
		return
	}

	results, err := (*engine).KNN(nodeUuid, int(kInt))
	if err != nil {
		writeError(err.Error(), w)
		return
	}

	json := simplejson.New()
	json.Set("uuid", nodeUuid.String())
	json.Set("results", results)

	response, err := json.MarshalJSON()
	if err != nil {
		writeError(err.Error(), w)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}
