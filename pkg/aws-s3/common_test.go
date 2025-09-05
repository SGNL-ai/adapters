// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
	"fmt"
	"strings"

	framework "github.com/sgnl-ai/adapter-framework"
	s3_adapter "github.com/sgnl-ai/adapters/pkg/aws-s3"
	cmnConfig "github.com/sgnl-ai/adapters/pkg/config"
	"github.com/sgnl-ai/adapters/pkg/testutil"
)

var (
	validAuthCredentials = &framework.DatasourceAuthCredentials{
		Basic: &framework.BasicAuthCredentials{
			Username: "test-username",
			Password: "test-password",
		},
	}

	validCommonConfig = &s3_adapter.Config{
		CommonConfig: &cmnConfig.CommonConfig{
			RequestTimeoutSeconds: testutil.GenPtr(120),
		},
		Region: "us-west-1",
		Bucket: "test-adapter-bucket",
	}

	// nolint: lll
	validCSVData = `Score,Customer Id,First Name,Last Name,Company,City,Country,Phone 1,Phone 2,Email,Subscription Date,Website,KnownAliases
1.1,e685B8690f9fbce,Erik,Little,Blankenship PLC,Caitlynmouth,Sao Tome and Principe,457-542-6899,055.415.2664x5425,shanehester@campbell.org,2021-12-23,https://wagner.com/,"[{""alias"": ""Shane Hester"", ""primary"": true},{""alias"": ""Cheyne Hester"", ""primary"": false}]"
2.2,6EDdBA3a2DFA7De,Yvonne,Shaw,Jensen and Sons,Janetfort,Palestinian Territory,9610730173,531-482-3000x7085,kleinluis@vang.com,2021-01-01,https://www.paul.org/,"[{""primary"": true, ""alias"": ""Klein Luis""},{""alias"": ""Cline Luis"", ""primary"": false}]"
3.3,b9Da13bedEc47de,Jeffery,Ibarra,"Rose, Deleon and Sanders",Darlenebury,Albania,(840)539-1797x479,209-519-5817,deckerjamie@bartlett.biz,2020-03-30,https://www.morgan-phelps.com/,"[{""alias"": ""Decker Jaime"", ""primary"": true}"
4.4,710D4dA2FAa96B5,James,Walters,Kline and Sons,Donhaven,Bahrain,+1-985-596-1072x3040,(528)734-8924x054,dochoa@carey-morse.com,2022-01-18,https://brennan.com/,"[{""alias"": ""Do Choa"", ""primary"": true}]"
5.5,3c44ed62d7BfEBC,Leslie,Snyder,"Price, Mason and Doyle",Mossfort,Central African Republic,812-016-9904x8231,254.631.9380,darrylbarber@warren.org,2020-01-25,http://www.trujillo-sullivan.info/,"[{""alias"": ""Darryl Barber"", ""primary"": true}"`

	// CSV data with integer scores for Int64 conversion testing
	// nolint: lll
	integerCSVData = `Score,Customer Id,First Name,Last Name,Company,City,Country,Phone 1,Phone 2,Email,Subscription Date,Website,KnownAliases
11,e685B8690f9fbce,Erik,Little,Blankenship PLC,Caitlynmouth,Sao Tome and Principe,457-542-6899,055.415.2664x5425,shanehester@campbell.org,2021-12-23,https://wagner.com/,"[{""alias"": ""Shane Hester"", ""primary"": true},{""alias"": ""Cheyne Hester"", ""primary"": false}]"
22,6EDdBA3a2DFA7De,Yvonne,Shaw,Jensen and Sons,Janetfort,Palestinian Territory,9610730173,531-482-3000x7085,kleinluis@vang.com,2021-01-01,https://www.paul.org/,"[{""primary"": true, ""alias"": ""Klein Luis""},{""alias"": ""Cline Luis"", ""primary"": false}]"
33,b9Da13bedEc47de,Jeffery,Ibarra,"Rose, Deleon and Sanders",Darlenebury,Albania,(840)539-1797x479,209-519-5817,deckerjamie@bartlett.biz,2020-03-30,https://www.morgan-phelps.com/,"[{""alias"": ""Decker Jaime"", ""primary"": true}]"
44,710D4dA2FAa96B5,James,Walters,Kline and Sons,Donhaven,Bahrain,+1-985-596-1072x3040,(528)734-8924x054,dochoa@carey-morse.com,2022-01-18,https://brennan.com/,"[{""alias"": ""Do Choa"", ""primary"": true}]"
55,3c44ed62d7BfEBC,Leslie,Snyder,"Price, Mason and Doyle",Mossfort,Central African Republic,812-016-9904x8231,254.631.9380,darrylbarber@warren.org,2020-01-25,http://www.trujillo-sullivan.info/,"[{""alias"": ""Darryl Barber"", ""primary"": true}"`

	// CSV data with string-friendly numeric scores for String conversion testing
	// nolint: lll
	stringCSVData = `Score,Customer Id,First Name,Last Name,Company,City,Country,Phone 1,Phone 2,Email,Subscription Date,Website,KnownAliases
1.1,e685B8690f9fbce,Erik,Little,Blankenship PLC,Caitlynmouth,Sao Tome and Principe,457-542-6899,055.415.2664x5425,shanehester@campbell.org,2021-12-23,https://wagner.com/,"[{""alias"": ""Shane Hester"", ""primary"": true},{""alias"": ""Cheyne Hester"", ""primary"": false}]"
2.2,6EDdBA3a2DFA7De,Yvonne,Shaw,Jensen and Sons,Janetfort,Palestinian Territory,9610730173,531-482-3000x7085,kleinluis@vang.com,2021-01-01,https://www.paul.org/,"[{""primary"": true, ""alias"": ""Klein Luis""},{""alias"": ""Cline Luis"", ""primary"": false}]"
3.3,b9Da13bedEc47de,Jeffery,Ibarra,"Rose, Deleon and Sanders",Darlenebury,Albania,(840)539-1797x479,209-519-5817,deckerjamie@bartlett.biz,2020-03-30,https://www.morgan-phelps.com/,"[{""alias"": ""Decker Jaime"", ""primary"": true}"
4.4,710D4dA2FAa96B5,James,Walters,Kline and Sons,Donhaven,Bahrain,+1-985-596-1072x3040,(528)734-8924x054,dochoa@carey-morse.com,2022-01-18,https://brennan.com/,"[{""alias"": ""Do Choa"", ""primary"": true}]"
5.5,3c44ed62d7BfEBC,Leslie,Snyder,"Price, Mason and Doyle",Mossfort,Central African Republic,812-016-9904x8231,254.631.9380,darrylbarber@warren.org,2020-01-25,http://www.trujillo-sullivan.info/,"[{""alias"": ""Darryl Barber"", ""primary"": true}"`

	// CSV data with decimal scores for Double conversion testing  
	// nolint: lll
	doubleCSVData = `Score,Customer Id,First Name,Last Name,Company,City,Country,Phone 1,Phone 2,Email,Subscription Date,Website,KnownAliases
1.1,e685B8690f9fbce,Erik,Little,Blankenship PLC,Caitlynmouth,Sao Tome and Principe,457-542-6899,055.415.2664x5425,shanehester@campbell.org,2021-12-23,https://wagner.com/,"[{""alias"": ""Shane Hester"", ""primary"": true},{""alias"": ""Cheyne Hester"", ""primary"": false}]"
2.2,6EDdBA3a2DFA7De,Yvonne,Shaw,Jensen and Sons,Janetfort,Palestinian Territory,9610730173,531-482-3000x7085,kleinluis@vang.com,2021-01-01,https://www.paul.org/,"[{""primary"": true, ""alias"": ""Klein Luis""},{""alias"": ""Cline Luis"", ""primary"": false}]"
3.3,b9Da13bedEc47de,Jeffery,Ibarra,"Rose, Deleon and Sanders",Darlenebury,Albania,(840)539-1797x479,209-519-5817,deckerjamie@bartlett.biz,2020-03-30,https://www.morgan-phelps.com/,"[{""alias"": ""Decker Jaime"", ""primary"": true}"
4.4,710D4dA2FAa96B5,James,Walters,Kline and Sons,Donhaven,Bahrain,+1-985-596-1072x3040,(528)734-8924x054,dochoa@carey-morse.com,2022-01-18,https://brennan.com/,"[{""alias"": ""Do Choa"", ""primary"": true}]"
5.5,3c44ed62d7BfEBC,Leslie,Snyder,"Price, Mason and Doyle",Mossfort,Central African Republic,812-016-9904x8231,254.631.9380,darrylbarber@warren.org,2020-01-25,http://www.trujillo-sullivan.info/,"[{""alias"": ""Darryl Barber"", ""primary"": true}"`

	// nolint: lll
	headersOnlyCSVData = `Score,Customer Id,First Name,Last Name,Company,City,Country,Phone 1,Phone 2,Email,Subscription Date,Website`

	// This is a TSV.
	// nolint: lll
	corruptCSVData = `Score	Customer Id	First Name	Last Name	Company	City	Country	Phone 1	Phone 2	Email	Subscription Date	Website
1	e685B8690f9fbce	Erik	Little	Blankenship PLC	Caitlynmouth	Sao Tome and Principe	457-542-6899	055.415.2664x5425	shanehester@campbell.org	2021-12-23	https://wagner.com/
2	6EDdBA3a2DFA7De	Yvonne	Shaw	Jensen and Sons	Janetfort	Palestinian Territory	9610730173	531-482-3000x7085	kleinluis@vang.com	2021-01-01	https://www.paul.org/
3	b9Da13bedEc47de	Jeffery	Ibarra	"Rose	 Deleon and Sanders"	Darlenebury	Albania	(840)539-1797x479	209-519-5817	deckerjamie@bartlett.biz	2020-03-30	https://www.morgan-phelps.com/
4	710D4dA2FAa96B5	James	Walters	Kline and Sons	Donhaven	Bahrain	+1-985-596-1072x3040	(528)734-8924x054	dochoa@carey-morse.com	2022-01-18	https://brennan.com/
5	3c44ed62d7BfEBC	Leslie	Snyder	"Price	 Mason and Doyle"	Mossfort	Central African Republic	812-016-9904x8231	254.631.9380	darrylbarber@warren.org	2020-01-25	http://www.trujillo-sullivan.info/`
)

// generateLargeCSVData creates a large CSV string.
func generateLargeCSVData() string {
	header := "Score,Customer Id,First Name,Last Name,Company,City,Country," +
		"Phone 1,Phone 2,Email,Subscription Date,Website,KnownAliases\n"

	const targetRows = 50000

	var builder strings.Builder

	builder.WriteString(header)

	builder.Grow(10*1024*1024 + 1024*1024) // 11MiB capacity (does not have any significance, can be any big number)

	for i := 1; i <= targetRows; i++ {
		score := float64(i) * 0.1
		customerID := fmt.Sprintf("ID%07d", i)
		firstName := fmt.Sprintf("FirstName%d", i)
		lastName := fmt.Sprintf("LastName%d", i)
		company := fmt.Sprintf("Company %d LLC", i)
		city := fmt.Sprintf("City%d", i)
		country := "TestCountry"
		phone1 := fmt.Sprintf("555-%03d-%04d", i%1000, i%10000)
		phone2 := fmt.Sprintf("666-%03d-%04d", (i+1)%1000, (i+1)%10000)
		email := fmt.Sprintf("user%d@example.com", i)
		subDate := "2024-01-01"
		website := fmt.Sprintf("https://example%d.com", i)
		aliases := fmt.Sprintf(`"[{""alias"": ""Alias%d"", ""primary"": true}]"`, i)

		row := fmt.Sprintf("%.1f,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
			score, customerID, firstName, lastName, company, city, country,
			phone1, phone2, email, subDate, website, aliases)

		builder.WriteString(row)
	}

	return builder.String()
}
