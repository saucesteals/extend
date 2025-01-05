# Extend API

Programmatically manage virtual cards on [Extend](https://www.paywithextend.com/)

> [!CAUTION]
> This is an **unofficial** client that is not affiliated with or endorsed by Extend.
> Using this client may violate Extend's Terms of Service.
> By using this client, you acknowledge and accept all associated risks.

## Installation

```sh
go get github.com/saucesteals/extend
```

## Device Verification

1. Login to the Extend dashboard
2. Open the browser's developer tools
3. Run the following JavaScript in the console:

```js
function lookupUserDetails() {
  const lastAuthUser = Object.entries(localStorage).find(
    ([key]) =>
      key.startsWith("CognitoIdentityServiceProvider") &&
      key.endsWith("LastAuthUser")
  );
  if (!lastAuthUser) return console.log("No user logged in");

  const [key, user] = lastAuthUser;
  const prefix = key.split(".").slice(0, 2).join(".") + "." + user + ".";
  const email = JSON.parse(
    localStorage[prefix + "userData"]
  ).UserAttributes.find((attr) => attr.Name === "email").Value;
  const lookup = (suffix) =>
    console.log(suffix, "-", localStorage[prefix + suffix]);

  console.log("Details for", email);
  lookup("deviceGroupKey");
  lookup("deviceKey");
  lookup("randomPasswordKey");
}

lookupUserDetails();
```

## Quick Start

### Initialize the client

```go
client := extend.New(cognito.NewCognito(cognito.AuthParams{
	Username:       "user@email.com",
	Password:       "password",

	DeviceGroupKey: "device_group_key", // deviceGroupKey from browser
	DeviceKey:      "device_key", // deviceKey from browser
	DevicePassword: "device_password", // randomPasswordKey from browser
}))
```

### Create a virtual card

```go
card, err := client.CreateVirtualCard(extend.CreateVirtualCardOptions{
	CreditCardID: "cc_id",
	DisplayName:  "Team Expenses",
	BalanceCents: 10000,
	Currency:     extend.CurrencyUSD,
	ValidTo:      time.Now().AddDate(0, 1, 0),
	Recipient:    "team@company.com",
	Notes:        "This card is for team expenses",
})
```

### Get a virtual card

```go
card, err := client.GetVirtualCard("vc_id")
```

### Cancel a virtual card

```go
card, err := client.CancelVirtualCard("vc_id")
```

### Close a virtual card

```go
card, err := client.CloseVirtualCard("vc_id")
```

### List virtual cards with pagination

```go
cards := client.ListVirtualCards(&extend.ListVirtualCardsOptions{
    PaginationOptions: extend.PaginationOptions{
		Page:          0,
		Count:         10,
		SortDirection: extend.SortDirectionAsc,
		SortField:     "activeClosedUpdatedAt",
	},
	CardholderOrViewer: "me",
	Issued:             true,
	Statuses:           []extend.VirtualCardStatus{extend.VirtualCardStatusActive},
})

for cards.Next() {
	page, err := cards.Get()
	for _, card := range page.Items() {
		// Process each card
	}
}
```

### Bulk create virtual cards

```go
cards := []extend.BulkCreateVirtualCard{
	{
		CardType:     extend.VirtualCardTypeStandard,
		Recipient:    "user1@company.com",
		DisplayName:  "Marketing Card 1",
		BalanceCents: 10000,
		ValidTo:      time.Now().AddDate(0, 1, 0),
	},
	{
		CardType:     extend.VirtualCardTypeStandard,
		Recipient:    "user2@company.com",
		DisplayName:  "Marketing Card 2",
		BalanceCents: 20000,
		ValidTo:      time.Now().AddDate(0, 1, 0),
	},
}

upload, err := client.BulkCreateVirtualCards("cc_id", cards)

// Check bulk upload status
status, err := client.GetBulkVirtualCardUpload(upload.BulkVirtualCardPush.BulkVirtualCardUploadID)
```
