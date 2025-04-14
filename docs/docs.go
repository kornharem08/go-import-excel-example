// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@example.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/purchaseorders": {
            "post": {
                "description": "Retrieves purchase order data from an Excel file located on a fixed network share path",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "purchaseorders"
                ],
                "summary": "Import purchase orders from Excel file on network share",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter orders by Job ID No",
                        "name": "job_id_no",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Path to the Excel file",
                        "name": "path",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "array",
                                "items": {
                                    "$ref": "#/definitions/models.PurchaseOrder"
                                }
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/purchaseorders/setting": {
            "get": {
                "description": "Retrieves the path of the purchase order Excel file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "purchaseorders"
                ],
                "summary": "Get the path of the purchase order Excel file",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.PurchaseOrder": {
            "type": "object",
            "properties": {
                "customer": {
                    "type": "string"
                },
                "delivery_date": {
                    "type": "string"
                },
                "distribution": {
                    "type": "string"
                },
                "job_id_no": {
                    "type": "string"
                },
                "ordered": {
                    "type": "integer"
                },
                "payment_term": {
                    "type": "string"
                },
                "po": {
                    "type": "string"
                },
                "po_date": {
                    "type": "string"
                },
                "pr": {
                    "type": "string"
                },
                "pr_date": {
                    "type": "string"
                },
                "product_code": {
                    "type": "string"
                },
                "product_description": {
                    "type": "string"
                },
                "project_manager": {
                    "type": "string"
                },
                "received": {
                    "type": "integer"
                },
                "remain": {
                    "type": "integer"
                },
                "request_date": {
                    "type": "string"
                },
                "sales_team": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Excel Order API",
	Description:      "Type \"Bearer\" followed by a space and JWT token.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
