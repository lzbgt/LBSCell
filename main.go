// radio,mcc,net,area,cell,unit,lon,lat,range,samples,changeable,created,updated,averageSignal
// GSM,460,0,4282,1374,,116.435411287,39.9098614727,493,55,1,1377178208,1432976121,
// lon, lat
// 1,2,3,4, 6,7
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

type Location struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

var Locations map[string]*Location = make(map[string]*Location)

func buildHash(params ...string) string {
	return strings.Join(params, ":")
}

func main() {
	// load MLS database
	flagPort := flag.String("port", "8081", "port")
	flagPath := flag.String("path", "MLS-460.csv", "csv path")
	loadMLS(*flagPath)
	// start the embedded web server
	r := mux.NewRouter()
	r.HandleFunc("/api/{component}", handler)
	http.Handle("/", r)
	http.ListenAndServe("0.0.0.0:"+*flagPort, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	var ret []byte = nil
	var err error = nil
	vars := mux.Vars(r)
	coapi := vars["component"]
	switch coapi {
	case "lbs":
		mcc := r.FormValue("mcc")
		mnc := r.FormValue("mnc")
		lac := r.FormValue("lac")
		cell := r.FormValue("cell")
		loc, ok := Locations[buildHash(mcc, mnc, lac, cell)]
		if ok && loc != nil {
			ret, err = json.Marshal(loc)
			if err != nil {
				//
				return
			}
		}

	default:
		ret = []byte("{\"success\":false, \"msg\":\"unknown api\"}")

	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func loadMLS(path string) {
	csvfile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1
	// see the Reader struct information below
	for record, err := reader.Read(); err == nil; record, err = reader.Read() {
		if len(record) > 7 {
			Locations[buildHash(record[1], record[2], record[3], record[4])] =
				&Location{record[7], record[6]}
		}
	}
}
