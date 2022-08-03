/*
 *  Copyright © 2016 Iury Braun
 *  Copyright © 2017 Weyboo
 */

package sdtp

import (
    //"github.com/vmihailenco/msgpack"
    //"weyboo.com/wws/lib/sdtp/response"
)

type R1 map[string]interface{}
type E1 []map[string]interface{}


func New() R1 {
    return make(R1)
}

func NewErrorData() E1 {
    return make(E1, 0)
}


func AddId(x R1, value int) {
    x["id"] = value
}

func AddError(x R1, error interface{}) {
    x["error"] = error
}

func AddErrorData(x E1, error map[string]interface{}) E1 {
    return append(x, error)
}

func AddResult(x R1, value interface{}) {
    x["result"] = value
}
