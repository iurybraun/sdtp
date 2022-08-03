/*
 * Copyright © 2016 Iury Braun
 * Copyright © 2017 Weyboo
 */

package sdtp

const (
    /* Standard errors */
  	ParseError     = -32700     /* Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text. */
    InvalidRequest = -32600     /* The JSON sent is not a valid Request object. */
    MethodNotFound = -32601     /* The method does not exist / is not available. */
    InvalidParams  = -32602     /* Invalid method parameter(s). */
    InternalError  = -32603     /* Internal JSON-RPC error. */
    
    /* Security errors */
    Unauthorized = -401
    Forbidden    = -403
    
    /* Resource errors */
    NotFound            = -404
    
    /* Process errors */
    Conflict            = -409
    UnprocessableEntity = -422
    
    
    /* Parameter errors */
    Required       = 1     /* */
    TooLow         = 2     /* Should be used when a to low value of Field was given. */
    LimitExceeded  = 3     /* Should be used when a limit is exceeded, e.g. for the field limit in a block. */
    Rejected       = 4     /* Should be used when an action was rejected, e.g. because of its content (too long contract code, containing wrong characters ?, should differ from -32602 - Invalid params). */
)
/**
    -	        200 Ok                          200     {"id": <ID>, "result": <DATA>}
    -	        201 Created                     200     {"id": <ID>, "result": <DATA>}
    security	401 Unauthorized                200     {"id": <ID>, "error": {"code": -401, "message": <MESSAGE>, "data": []}}
    security	403 Forbidden                   200     {"id": <ID>, "error": {"code": -403, "message": <MESSAGE>, "data": []}}
    -	        404 Not Found                   200     {"id": <ID>, "result": null}
    validation	409 Conflict                    200     {"id": <ID>, "error": {"code": -409, "message": <MESSAGE>, "data": []}}
    execution   422 Unprocessable Entity        200     {"id": <ID>, "error": {"code": -422, "message": <MESSAGE>}}
    system	    500 Internal Server Error       200     {"id": <ID>, "error": {"code": -32603, "message": <MESSAGE>}}
    
    
    
    PROCESS TO RESPONSE
    -> Standard errors
    -> Security errors
    -> Process errors
    
    
    CREATE                     409  {"id": <ID>, "result": <DATA>}
    READ  *                         {"id": <ID>, "result": {"meta": <META>, "data": <DATA>, "pagination": <PAGINATION>}}
          obj              404      {"id": <ID>, "result": <DATA>}
    UPDATE/DELETE          404 409  {"id": <ID>, "result": true}
*/
/**
   Standard Errors
     ParseError
     InvalidRequest
     MethodNotFound
     InvalidParams
        < Parameter errors :: Required, TooLow, LimitExceeded, Rejected >       -       field, file
     InternalError
   
   Security Errors
     Unauthorized < empty and string >
        < Parameter errors :: Required, TooLow, LimitExceeded, Rejected >       -       field, file
     Forbidden
   
   Resource Errors
     NotFound
   
   Process Errors
     Conflict < empty and string >
        < Parameter errors :: Rejected >       -       field, file
     UnprocessableEntity
*/
