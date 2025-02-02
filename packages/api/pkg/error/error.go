package error

import "errors"

var ErrGeneratingToken = errors.New("error generating token")
var ErrCreatingToken = errors.New("error creating token in database")

// assets
var ErrCreatingAsset = errors.New("error creating asset in database")
var ErrGettingAsset = errors.New("error getting asset in database")
var ErrGettingAssets = errors.New("error getting assets in database")
var ErrUpdatingAsset = errors.New("error updating asset in database")
var ErrDeletingAsset = errors.New("error deleting asset in database")

var ErrCreatingAssetIDExists = errors.New("id exists")
var ErrUploadingFile = errors.New("error uploading file")
var ErrAssetDoesNotBelongToTheUser = errors.New("asset does not exist or do not belong to the user")
