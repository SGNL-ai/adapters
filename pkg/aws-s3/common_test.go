// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst

package awss3_test

import (
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
