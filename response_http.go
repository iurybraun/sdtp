/*
 * Copyright © 2016 Iury Braun
 * Copyright © 2017 Weyboo
 * 
 * Return:
 *  Content-type :: "application/json; charset=utf-8", "application/msgpack" (default)
 * 
 * Request:
 *  Content-type :: "application/json;", "application/msgpack"
 *  Accept :: "application/json"
 */

package sdtp

import (
    "time"
    "net/http"
    "encoding/json"
    
	"github.com/vmihailenco/msgpack"
    
    "github.com/iurybraun/go-cfg_ini"
    "github.com/iurybraun/i18n_ini"
)

func SendHttpData(data interface{}, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    if data == nil {
        SendHttpInternalError(start_time, w, r)
        return
    }
    
    
    result_map := New()
    
    AddResult(result_map, data)
    
    if r.Header.Get("Accept") == "application/json" {
        msgjson, err := json.Marshal(result_map)
        if err != nil {
            SendHttpInternalError(start_time, w, r)
            return
        }
        
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.Write(msgjson)
        
        //logsHttpCreate(nil, start_time, int32(len(msgjson)), w, r)
    } else {
        msgpk, err := msgpack.Marshal(result_map)
        if err != nil {
            SendHttpInternalError(start_time, w, r)
            return
        }
        
        w.Header().Set("Content-Type", "application/msgpack")
        w.Write(msgpk)
        
        //logsHttpCreate(nil, start_time, int32(len(msgpk)), w, r)
    }
}

func sendHttpErr(err_d interface{}, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    result_map := New()
    
    AddError(result_map, err_d)
    
    if r.Header.Get("Accept") == "application/json" {
        msgjson, err := json.Marshal(result_map)
        if err != nil {
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            w.Write(nil)
            return
        }
        
        w.Header().Set("Content-Type", "application/json; charset=utf-8")
        w.Write(msgjson)
        
        //logsHttpCreate(err_d, start_time, int32(len(msgjson)), w, r)
    } else {
        msgpk, err := msgpack.Marshal(result_map)
        if err != nil {
            w.Header().Set("Content-Type", "application/msgpack")
            w.Write(nil)
            return
        }
        
        w.Header().Set("Content-Type", "application/msgpack")
        w.Write(msgpk)
        
        //logsHttpCreate(err_d, start_time, int32(len(msgpk)), w, r)
    }
}



/**
 * Standard errors
 */
func SendHttpParseError(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": ParseError,
            "message": i18n_ini.LoadTr(lang, "standard-errors.ParseError"),
        }, start_time, w, r)
}

func SendHttpInvalidRequest(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": InvalidRequest,
            "message": i18n_ini.LoadTr(lang, "standard-errors.InvalidRequest"),
        }, start_time, w, r)
}

func SendHttpMethodNotFound(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": MethodNotFound,
            "message": i18n_ini.LoadTr(lang, "standard-errors.MethodNotFound"),
        }, start_time, w, r)
}

func SendHttpSingleInvalidParams(parameter_error int, location string, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    var message_parameter_error string
    
    switch parameter_error {
        case 1:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.Required")
        break;
        case 2:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.TooLow")
        break;
        case 3:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.LimitExceeded")
        break;
        case 4:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.Rejected")
        break;
    }
    
    
    err_map := NewErrorData()
    err_map = AddErrorData(err_map, map[string]interface{}{
            "code": parameter_error,
            "location": location,
            "message": message_parameter_error,
        })

    sendHttpErr(
        map[string]interface{}{
            "code": InvalidParams,
            "message": i18n_ini.LoadTr(lang, "standard-errors.InvalidParams"),
            "data": err_map,
        }, start_time, w, r)
}

/*
    err_map := NewErrorData()
    err_map = AddErrorData(err_map, map[string]interface{}{
            "code": parameter_error,
            "location": location,
            "message": i18n_ini.LoadTr(locale, "parameter-errors.Rejected"),
        })
    
    SendHttpMultipleInvalidParams(err_map, w, r)
*/
func SendHttpMultipleInvalidParams(err_map E1, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": InvalidParams,
            "message": i18n_ini.LoadTr(lang, "standard-errors.InvalidParams"),
            "data": err_map,
        }, start_time, w, r)
}

func SendHttpInternalError(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": InternalError,
            "message": i18n_ini.LoadTr(lang, "standard-errors.InternalError"),
        }, start_time, w, r)
}


/**
 * Security errors
 */
func SendHttpUnauthorized(message string, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    var message_t string
    
    if message != "" {
        message_t = message
    } else {
        message_t = i18n_ini.LoadTr(lang, "security-errors.Unauthorized")
    }
    
    sendHttpErr(
        map[string]interface{}{
            "code": Unauthorized,
            "message": message_t,
        }, start_time, w, r)
}

func SendHttpSingleUnauthorized(parameter_error int, location string, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    var message_parameter_error string
    
    switch parameter_error {
        case 1:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.Required")
        break;
        case 2:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.TooLow")
        break;
        case 3:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.LimitExceeded")
        break;
        case 4:
            message_parameter_error = i18n_ini.LoadTr(lang, "parameter-errors.Rejected")
        break;
    }
    
    
    err_map := NewErrorData()
    err_map = AddErrorData(err_map, map[string]interface{}{
            "code": parameter_error,
            "location": location,
            "message": message_parameter_error,
        })

    sendHttpErr(
        map[string]interface{}{
            "code": Unauthorized,
            "message": i18n_ini.LoadTr(lang, "security-errors.Unauthorized"),
            "data": err_map,
        }, start_time, w, r)
}

func SendHttpForbidden(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": Forbidden,
            "message": i18n_ini.LoadTr(lang, "security-errors.Forbidden"),
        }, start_time, w, r)
}



/**
 * Resource errors
 */
func SendHttpNotFound(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": NotFound,
            "message": i18n_ini.LoadTr(lang, "resource-errors.NotFound"),
        }, start_time, w, r)
}



/**
 * Process errors
 */
func SendHttpConflict(message string, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    var message_t string
    
    if message != "" {
        message_t = message
    } else {
        message_t = i18n_ini.LoadTr(lang, "process-errors.Conflict")
    }
    
    sendHttpErr(
        map[string]interface{}{
            "code": Conflict,
            "message": message_t,
        }, start_time, w, r)
}

func SendHttpSingleConflict(location string, start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    err_map := NewErrorData()
    err_map = AddErrorData(err_map, map[string]interface{}{
            "code": 4,
            "location": location,
            "message": i18n_ini.LoadTr(lang, "parameter-errors.Rejected"),
        })

    sendHttpErr(
        map[string]interface{}{
            "code": Conflict,
            "message": i18n_ini.LoadTr(lang, "process-errors.Conflict"),
            "data": err_map,
        }, start_time, w, r)
}

func SendHttpUnprocessableEntity(start_time time.Time, w http.ResponseWriter, r *http.Request) {
    var lang = cfg_ini.GetLang(r)
    
    sendHttpErr(
        map[string]interface{}{
            "code": UnprocessableEntity,
            "message": i18n_ini.LoadTr(lang, "process-errors.UnprocessableEntity"),
        }, start_time, w, r)
}


/*
err_map := sdtp.NewErrorData()
err_map = sdtp.AddErrorData(err_map, map[string]interface{}{
        "code": -2,
        "message": i18n_ini.LoadTr(locale, "crud.EmptyValue"),
        "location_type": "param2",
        "location": "state2",
    })

sdtp.SendHttpErr(
    map[string]interface{}{
        "code": sdtp.MethodNotFound,
        "message": i18n_ini.LoadTr(locale, "standard-errors.MethodNotFound1"),
        "data": err_map,
    }, w, r)
*/
