package phonenumbers

import (
	"fmt"
	"log"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)

type PhoneNumberProvider struct {
	client *twilio.RestClient
}

func NewPhoneNumberProvider() PhoneNumberProvider {
	client := twilio.NewRestClient()

	return PhoneNumberProvider{client: client}
}

func (p PhoneNumberProvider) getAvailablePhoneNumber() *api.ApiV2010AvailablePhoneNumberLocal {
	var phoneNumber string
	params := &api.ListAvailablePhoneNumberLocalParams{}
	params.SetAreaCode(415)
	params.SetLimit(1)
	params.SetExcludeAllAddressRequired(true)
	params.SetSmsEnabled(true)

	resp, err := p.client.Api.ListAvailablePhoneNumberLocal("US", params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for record := range resp {
			if resp[record].FriendlyName != nil {
				phoneNumber = *resp[record].PhoneNumber
				log.Printf("New phone number found, going to register: %v", phoneNumber)
				return &resp[record]
			}
		}
	}
	return nil
}

func (p PhoneNumberProvider) registerPhoneNumber(phoneNumber string) *string {
	params := &api.CreateIncomingPhoneNumberParams{}
	params.SetPhoneNumber(phoneNumber)
	params.SetSmsApplicationSid("AP16d4fd44166f87d307778dec512cd1cc")

	resp, err := p.client.Api.CreateIncomingPhoneNumber(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Status != nil {
			return resp.Status
		}
	}
	return nil
}

func (p PhoneNumberProvider) ProvideNewNumber() string {
	resp := p.getAvailablePhoneNumber()
	if resp == nil {
		panic("Could not find available phone number")
	}
	status := p.registerPhoneNumber(*resp.FriendlyName)
	if status == nil {
		panic("Could not register phone number")
	}
	log.Printf("Phone number %s registered with status %s\n", *resp.PhoneNumber, *status)
	return *resp.PhoneNumber
}
