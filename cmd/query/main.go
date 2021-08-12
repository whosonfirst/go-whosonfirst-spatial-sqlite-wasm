package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-sqlite"
	"github.com/whosonfirst/go-whosonfirst-spatial/api"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"log"
	"syscall/js"	
)

func pip(req *api.PointInPolygonRequest) (spr.StandardPlacesResponse, error) {

	c, err := geo.NewCoordinate(req.Longitude, req.Latitude)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to create new coordinate, %v", err)
	}
	
	f, err := api.NewSPRFilterFromPointInPolygonRequest(req)
	
	if err != nil {
		return nil, err
	}
	
	r, err := db.PointInPolygon(ctx, c, f)
	
	if err != nil {
		return nil, fmt.Errorf("Failed to query database with coord %v, %v", c, err)
	}
	
	return r, nil
}

func queryFunc(this js.Value, args []js.Value) interface{} {

	str_req := args[0].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		resolve := args[0]
		reject := args[1]

		var req *api.PointInPolygonRequest

		err := json.Unmarshal([]byte(str_req), &req)
			
		if err != nil {
			reject.Invoke(fmt.Sprintf("Failed to parse request, %v", err))
			return nil
		}
		
		rsp, err := pip(req)

		if err != nil {
			reject.Invoke(fmt.Sprintf("Failed to query, %v", err))
			return nil
		}

		enc_rsp, err := json.Marshal(rsp)

		if err != nil {
			reject.Invoke(fmt.Sprintf("Failed to marshal response, %v", err))
			return nil
		}

		resolve.Invoke(string(enc_rsp))
		return nil
	}
		
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)		
}
	
func main() {

	fs, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	flagset.Parse(fs)

	err = flags.ValidateCommonFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	database_uri, _ := lookup.StringVar(fs, flags.SPATIAL_DATABASE_URI)

	ctx := context.Background()
	db, err := database.NewSpatialDatabase(ctx, database_uri)

	if err != nil {
		log.Fatalf("Failed to create database for '%s', %v", database_uri, err)
	}

	query_func := queryeFunc(db)
	defer query_func.Release()

	js.Global().Set("point_in_point", query_func)

	c := make(chan struct{}, 0)

	log.Println("WASM point-in-polygon function initialized")
	<-c
}
