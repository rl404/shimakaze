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
        "/agencies": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agency"
                ],
                "summary": "Get agency data.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/service.agency"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/agencies/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Agency"
                ],
                "summary": "Get agency data.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "wikia id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/service.agency"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber data.",
                "parameters": [
                    {
                        "enum": [
                            "all",
                            "stats"
                        ],
                        "type": "string",
                        "default": "all",
                        "description": "mode",
                        "name": "mode",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "names",
                        "name": "names",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "original name",
                        "name": "original_name",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "nickname",
                        "name": "nickname",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "exclude active",
                        "name": "exclude_active",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "exclude retired",
                        "name": "exclude_retired",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start debut year",
                        "name": "start_debut_year",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "end debut year",
                        "name": "end_debut_year",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "start retired year",
                        "name": "start_retired_year",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "end retired year",
                        "name": "end_retired_year",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "has 2d model",
                        "name": "has_2d",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "has 3d model",
                        "name": "has_3d",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "character designer",
                        "name": "character_designer",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "character 2d modeler",
                        "name": "character_2d_modeler",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "character 3d modeler",
                        "name": "character_3d_modeler",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "in agency",
                        "name": "in_agency",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "agency",
                        "name": "agency",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "channel types",
                        "name": "channel_types",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "name",
                            "-name",
                            "debut_date",
                            "-debut_date",
                            "retirement_date",
                            "-retirement_date"
                        ],
                        "type": "string",
                        "default": "name",
                        "description": "sort",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "page",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 20,
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/service.vtuber"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/2d-modelers": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber character 2D modelers.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/3d-modelers": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber character 3D modelers.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/agency-trees": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber agency trees.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/service.vtuberAgencyTree"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/character-designers": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber character designers.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "type": "string"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/family-trees": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber family trees.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/service.vtuberFamilyTree"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/images": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get all vtuber images.",
                "parameters": [
                    {
                        "type": "boolean",
                        "description": "shuffle",
                        "name": "shuffle",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "limit",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "type": "array",
                                            "items": {
                                                "$ref": "#/definitions/service.vtuberImage"
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/vtubers/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Vtuber"
                ],
                "summary": "Get vtuber data.",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "wikia id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/utils.Response"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/service.vtuber"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        },
        "/wikia/image/{path}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Wikia"
                ],
                "summary": "Get wikia image.",
                "parameters": [
                    {
                        "type": "string",
                        "description": "wikia image url",
                        "name": "path",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "PNG image"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/utils.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "service.agency": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "service.vtuber": {
            "type": "object",
            "properties": {
                "affiliations": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "age": {
                    "type": "number"
                },
                "agencies": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.vtuberAgency"
                    }
                },
                "birthday": {
                    "type": "string"
                },
                "blood_type": {
                    "type": "string"
                },
                "caption": {
                    "type": "string"
                },
                "channels": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.vtuberChannel"
                    }
                },
                "character_2d_modelers": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "character_3d_modelers": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "character_designers": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "debut_date": {
                    "type": "string"
                },
                "emoji": {
                    "type": "string"
                },
                "gender": {
                    "type": "string"
                },
                "has_2d": {
                    "type": "boolean"
                },
                "has_3d": {
                    "type": "boolean"
                },
                "height": {
                    "type": "number"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "nicknames": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "official_websites": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "original_names": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "retirement_date": {
                    "type": "string"
                },
                "social_medias": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "updated_at": {
                    "type": "string"
                },
                "weight": {
                    "type": "number"
                },
                "zodiac_sign": {
                    "type": "string"
                }
            }
        },
        "service.vtuberAgency": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "service.vtuberAgencyTree": {
            "type": "object",
            "properties": {
                "links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.vtuberAgencyTreeLink"
                    }
                },
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.vtuberAgencyTreeNode"
                    }
                }
            }
        },
        "service.vtuberAgencyTreeLink": {
            "type": "object",
            "properties": {
                "id1": {
                    "type": "integer"
                },
                "id2": {
                    "type": "integer"
                }
            }
        },
        "service.vtuberAgencyTreeNode": {
            "type": "object",
            "properties": {
                "agencies": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "has_retired": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "service.vtuberChannel": {
            "type": "object",
            "properties": {
                "type": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                }
            }
        },
        "service.vtuberFamilyTree": {
            "type": "object",
            "properties": {
                "links": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.vtuberFamilyTreeLink"
                    }
                },
                "nodes": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/service.vtuberFamilyTreeNode"
                    }
                }
            }
        },
        "service.vtuberFamilyTreeLink": {
            "type": "object",
            "properties": {
                "id1": {
                    "type": "integer"
                },
                "id2": {
                    "type": "integer"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "service.vtuberFamilyTreeNode": {
            "type": "object",
            "properties": {
                "has_retired": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "service.vtuberImage": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "utils.Response": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                },
                "meta": {
                    "type": "object"
                },
                "status": {
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Shimakaze API",
	Description:      "Shimakaze API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
