package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"unicode/utf8"

	"github.com/ggicci/httpin"
)

// Constants
const (
	MaxDescriptionLength = 1000
	MaxFileSize          = 5 * 1024 * 1024 // 5MB
	Base58Pattern        = `^[1-9A-HJ-NP-Za-km-z]+$`
	AlphanumericPattern  = `^[a-zA-Z0-9]+$`
)

// Validation helpers

func validateID(id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if !isBase58(id) {
		return errors.New("id must be a valid Base58 string")
	}
	return nil
}

func validateImage(image *string) error {
	if image != nil {
		_, err := url.ParseRequestURI(*image)
		if err != nil {
			return errors.New("image must be a valid URL")
		}
	}
	return nil
}

func validateDescription(description *string) error {
	if description != nil {
		if utf8.RuneCountInString(*description) > MaxDescriptionLength {
			return fmt.Errorf("description exceeds the maximum length of %d characters", MaxDescriptionLength)
		}
		if containsMaliciousContent(*description) {
			return errors.New("description contains malicious content")
		}
	}
	return nil
}

func validateSocial(social *map[string]string) error {
	if social != nil {
		if _, err := json.Marshal(social); err != nil {
			return errors.New("social must be a valid JSON object")
		}
		for key, urlStr := range *social {
			if _, err := url.ParseRequestURI(urlStr); err != nil {
				return fmt.Errorf("social URL for '%s' is invalid", key)
			}
		}
	}
	return nil
}

func isBase58(s string) bool {
	re := regexp.MustCompile(Base58Pattern)
	return re.MatchString(s)
}

func containsMaliciousContent(s string) bool {
	htmlRegex := regexp.MustCompile(`(?i)(<script|<iframe|<object|<embed|<form|<input|<img|<svg|<style|<link|<base|<meta|<frame)`)
	return htmlRegex.MatchString(s)
}

// Struct definitions

type AuthHeaders struct {
	Signature string `in:"header=x-signature"`
	Address   string `in:"header=x-address"`
	Message   string `in:"header=x-message"`
}

type NewAsset struct {
	ID          string             `json:"id"`
	Description *string            `json:"description"`
	Image       *string            `json:"image" validate:"nonzero"`
	Social      *map[string]string `json:"social"`
}

type UpdateAsset struct {
	Description *string            `json:"description"`
	Image       *string            `json:"image" validate:"nonzero"`
	Social      *map[string]string `json:"social"`
}

// Input structs with validations

type CreateAssetInput struct {
	*AuthHeaders
	NewAsset `in:"body=json;nonzero"`
}

func (c *CreateAssetInput) Validate() error {
	if err := validateID(c.ID); err != nil {
		return err
	}
	if err := validateImage(c.Image); err != nil {
		return err
	}
	if err := validateDescription(c.Description); err != nil {
		return err
	}
	if err := validateSocial(c.Social); err != nil {
		return err
	}
	return nil
}

type UpdateAssetInput struct {
	*AuthHeaders
	ID          string `in:"path=id"`
	UpdateAsset `in:"body=json;nonzero"`
}

func (c *UpdateAssetInput) Validate() error {
	if err := validateID(c.ID); err != nil {
		return err
	}
	if err := validateImage(c.Image); err != nil {
		return err
	}
	if err := validateDescription(c.Description); err != nil {
		return err
	}
	if err := validateSocial(c.Social); err != nil {
		return err
	}
	return nil
}

type DeleteAssetInput struct {
	*AuthHeaders
	ID string `in:"path=id"`
}

func (c *DeleteAssetInput) Validate() error {
	if err := validateID(c.ID); err != nil {
		return err
	}

	return nil
}

type GetAssetByIDInput struct {
	ID string `in:"path=id"`
}

func (c *GetAssetByIDInput) Validate() error {
	return validateID(c.ID)
}

type GetAssetsInput struct {
	Address *string `in:"query=address"`
	ID      *string `in:"query=id"`
	Order
	Pagination
}

func (c *GetAssetsInput) Validate() error {
	if c.ID != nil {
		if err := validateID(*c.ID); err != nil {
			return err
		}
	}
	if c.Address != nil {
		if len(*c.Address) != 48 {
			return errors.New("address is invalid")
		}
		re := regexp.MustCompile(AlphanumericPattern)
		if !re.MatchString(*c.Address) {
			return errors.New("address must only contain alphanumeric characters")
		}
	}
	if c.Order.Order != nil {
		validOrders := map[string]bool{"id": true, "address": true, "created_at": true}
		if !validOrders[*c.Order.Order] {
			return errors.New("order fields are: id, address, and created_at")
		}
	}
	c.Pagination.Validate()
	return nil
}

type UploadImageInput struct {
	ID   string       `in:"form=id"`
	File *httpin.File `in:"form=file"`
}

func (c *UploadImageInput) Validate() error {
	if err := validateID(c.ID); err != nil {
		return err
	}
	if c.File == nil {
		return errors.New("file is required")
	}
	if c.File.Size() > MaxFileSize {
		return fmt.Errorf("file size exceeds the maximum allowed limit of %dMB", MaxFileSize/(1024*1024))
	}
	return nil
}
