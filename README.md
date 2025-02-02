# Assets API

## Overview

The Assets API provides endpoints to manage assets. This document describes the available endpoints, their expected input, and output.

--- 

## Features

- Health check endpoint
- Authentication nonce retrieval
- Image upload
- Asset creation, retrieval, and management

--- 

## Authentication

The endpoints that require authentication needs to have three headers in their request:

- `"X-Message"`: A valid nonce obtained from the `GET /nonce` endpoint. This nonce is used to prevent replay attacks.
- `"X-Signature"`: The message signed with the private key of the account. This signature is used to verify the authenticity of the request.
- `"X-Address"`: The address of the account used to sign the message. This address is used to identify the user making the request.


## API Endpoints

### **GET /health**

#### Description
Checks if the API server is running.

#### Response
- **200 OK**: Indicates the server is operational.

#### Example Response
Empty

### **GET /nonce**

#### Description
Retrieves a nonce for authentication. It must be signed and set to the `MESSAGE` header in the endpoints that require authentication. 

#### Response
- **200 OK**: with the nonce.

#### Example Response

```json
{
  "nonce": "nonce"
}
```

### **GET /assets**

#### Description
Retrieves a list of assets.

### Query arguments
- **address**: string. It must be a valid polkadot address. 
- **id**: string
- **order**: string. It orders the results using a field. It must be `id`, `address` or `created_at`.
- **ascending**: bool. It defines if the order is ascending or descendingl It's used along with `order`.
- **limit**: int. The maximum number of items to return. If not specified or greater than the maximum limit, it defaults to 100.
- **offset**: int. The number of items to skip before starting to collect the result set. Defaults to 0 if not specified.

#### Response
- **200 OK** with the list of assets.
- **400 Bad Request** if the input is invalid.

#### Example Response

```json
{

    "ok": true,
    "data": [
        {
            "id": "asset_id",
            "description": "asset_description",
            "address": "owner",
            "image": "asset_image_url",
            "social": {
            "twitter": "twitter_handle",
            "facebook": "facebook_handle"
            }
        }
    ]
}
```

### **GET /assets/{id}**

#### Description
Retrieves an asset by its ID.

#### Response
- **200 OK** with the asset details.
- **404 Not Found** if the asset is not found.

#### Example Response
```json
{
  "ok": true,
  "data": {
    "id": "asset_id",
    "description": "asset_description",
    "image": "asset_image_url",
    "social": {
        "twitter": "twitter_handle",
        "facebook": "facebook_handle"
    }
  }
}
```
### **POST /assets**

#### Description
Creates a new asset. It requires authentication.

#### Request Body
```json
{
  "ok": true,
  "data": {
    "id": "asset_id",
    "description": "asset_description",
    "image": "asset_image_url",
    "social": {
        "twitter": "twitter_handle",
        "facebook": "facebook_handle"
    }
  }
}
```

#### Response
- **201 Created** with the created asset details.
- **400 Bad Request** if the input is invalid.
- **401 Unauthorized** if the authentication fails.

#### Example Response
```json
{
  "ok": true,
  "data": {
    "id": "asset_id",
    "description": "asset_description",
    "image": "asset_image_url",
    "social": {
        "twitter": "twitter_handle",
        "facebook": "facebook_handle"
    }
  }
}
```
### **PUT /assets/{id}**

#### Description
It updates an asset. Only its owner can do it. It requires authentication.

#### Request Body
```json
{
    "ok": true
}

```
Note: at least one field is required.

#### Response
- **200 OK** 
- **400 Bad Request** if the input is invalid.
- **401 Unauthorized** if the authentication fails.
- **404 Not Found** if the combination of `id` and `address` does not exist.

#### Example Response
```json
{
    "ok": true
}
```
### **DELETE /assets/{id}**

#### Description
It deletes an asset. Only its owner can do it. It requires authentication.

#### Response
- **200 OK** 
- **400 Bad Request** if the input is invalid.
- **401 Unauthorized** if the authentication fails.
- **404 Not Found** if the combination of `id` and `address` does not exist.

#### Example Response
```json
{
    "ok": true
}
```

## Prerequisites
- Go 1.16 or later
- A running instance of the database
- Environment variables configured 

## Testing
To run the tests, use the following command:

``` bash
go test ./...
```

## Contributing
Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (git checkout -b feature-branch).
3. Make your changes.
4. Commit your changes (git commit -am 'Add new feature').
5. Push to the branch (git push origin feature-branch).
6. Create a new Pull Request.