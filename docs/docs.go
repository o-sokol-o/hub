// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/lists": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get all lists",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Checklists"
                ],
                "summary": "Get All Checklists",
                "operationId": "get-all-lists",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler_api.ChecklistsResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "create checklist",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Checklists"
                ],
                "summary": "Create Checklist",
                "operationId": "create-list",
                "parameters": [
                    {
                        "description": "Checklist info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateChecklist"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler_api.idResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            }
        },
        "/api/lists/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get list by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Checklists"
                ],
                "summary": "Get Checklist By Id",
                "operationId": "get-list-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Checklist ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Checklist"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get update by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Checklists"
                ],
                "summary": "Update Checklist By Id",
                "operationId": "get-update-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Checklist ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Checklist info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateChecklist"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get delete by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Checklists"
                ],
                "summary": "Delete Checklist By Id",
                "operationId": "get-delete-by-id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Checklist ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            }
        },
        "/auth/sign-in": {
            "post": {
                "description": "login",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "SignIn",
                "operationId": "login",
                "parameters": [
                    {
                        "description": "credentials",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handler_api.signInInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "token",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            }
        },
        "/auth/sign-up": {
            "post": {
                "description": "create account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "SignUp",
                "operationId": "create-account",
                "parameters": [
                    {
                        "description": "account info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/domain.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            }
        },
        "/v1/sensor": {
            "get": {
                "description": "?api_key=aqen104Ur2zNX1Ykwv4:a39831d103eb4c0d \u0026f100=0.01\u0026f101=28\u0026f102=0\u0026f103=17.51\u0026f104=15.52\u0026f105=1072 \u0026f200=17.52\u0026f201=134.06\u0026f202=317\u0026f203=25.7000 \u0026f400=3.27\u0026f401=0.39\u0026f402=3.26\u0026f403=0.39\u0026f404=0.08\u0026f405=0.00 \u0026f4002=0.08 \u0026f11000=504\u0026f10001=24.31\u0026f10004=0",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "AquaHub"
                ],
                "summary": "AquaHub sensors data store",
                "operationId": "aquahub-sensor-data-store",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    },
                    "default": {
                        "description": "",
                        "schema": {
                            "$ref": "#/definitions/handler_api.statusResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Checklist": {
            "type": "object",
            "required": [
                "id",
                "title"
            ],
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Rocket Launch Description"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "active",
                        "disabled"
                    ],
                    "example": "active"
                },
                "title": {
                    "type": "string",
                    "example": "Rocket Launch"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "domain.UpdateChecklist": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Description Checklist"
                },
                "status": {
                    "type": "string",
                    "enum": [
                        "active",
                        "disabled"
                    ],
                    "example": "active"
                },
                "title": {
                    "type": "string",
                    "example": "Title Checklist"
                }
            }
        },
        "domain.User": {
            "type": "object",
            "required": [
                "email",
                "first_name",
                "last_name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "ae@ae.ae"
                },
                "first_name": {
                    "type": "string",
                    "example": "Andy"
                },
                "last_name": {
                    "type": "string",
                    "example": "Sokol"
                },
                "password": {
                    "type": "string",
                    "example": "jyWtbKg76by"
                }
            }
        },
        "handler_api.ChecklistsResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Checklist"
                    }
                }
            }
        },
        "handler_api.idResponse": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "handler_api.signInInput": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "handler_api.statusResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "Ok"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "AquaHub API",
	Description:      "API Server for AquaHub",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
