package handler

import (
	"accomodation-service/model"
	"accomodation-service/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccomodationHandler struct {
	Service *service.AccomodationService
}

func (handler *AccomodationHandler) CreateAccomodation(w http.ResponseWriter, r *http.Request) {
	var accomodation model.Accomodation
	err := json.NewDecoder(r.Body).Decode(&accomodation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = handler.Service.CreateAccomodation(&accomodation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

func (handler *AccomodationHandler) SearchAccomodation(w http.ResponseWriter, r *http.Request) {
	var accomodation model.AccomodationSearch
	err := json.NewDecoder(r.Body).Decode(&accomodation)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//pronadji sve na toj lokaciji
	/*err = handler.Service.SearchAccByLocation(&accomodation.Location)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}*/

}
