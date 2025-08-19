// Copyright 2025 SGNL.ai, Inc.

// nolint: goconst, lll
package crowdstrike_test

// This file documents the responses for the mock CrowdStrike server.

var (
	UserResponsePage1 = `{
    "data": {
        "entities": {
            "pageInfo": {
                "hasNextPage": true,
                "endCursor": "eyJyaXNrU2NvcmUiOjAuNjQ1NDg3MTMzOTk5OTk5OSwiX2lkIjoiNDVkYzQwZTItN2I3Yi00ZjM4LTlhYzctOThmNGEzNWIyNGUxIn0="
            },
            "nodes": [
                {
                    "archived": false,
                    "creationTime": "2024-05-15T15:29:10.000Z",
                    "earliestSeenTraffic": "2024-05-23T02:02:43.960Z",
                    "emailAddresses": [],
                    "entityId": "095b6929-44b9-4525-a0cc-9ef4552011f3",
                    "hasADDomainAdminRole": true,
                    "impactScore": 0.92,
                    "inactive": true,
                    "learned": true,
                    "markTime": null,
                    "mostRecentActivity": "2024-05-29T23:27:14.229Z",
                    "riskScore": 0.66,
                    "riskScoreSeverity": "MEDIUM",
                    "riskScoreWithoutLinkedAccounts": 0.6561,
                    "secondaryDisplayName": "CORP.SGNL.AI\\Wendolyn.Garber",
                    "shared": false,
                    "stale": true,
                    "watched": false,
                    "type": "USER",
                    "riskFactors": [
                        {
                            "score": 0.6,
                            "severity": "MEDIUM",
                            "type": "WEAK_PASSWORD_POLICY"
                        },
                        {
                            "score": 0.425,
                            "severity": "NORMAL",
                            "type": "STALE_ACCOUNT"
                        }
                    ],
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "Wendolyn Garber",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "cc1ea590-c660-450f-b35a-841d553fb32d",
                                "6b518e93-b160-47e7-b02d-34d41c9677d3"
                            ],
                            "creationTime": "2024-05-15T15:29:10.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": "Finance",
                            "description": null,
                            "dn": "CN=Wendolyn Garber,OU=Users,OU=Company,DC=corp,DC=sgnl,DC=ai",
                            "domain": "CORP.SGNL.AI",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "cc1ea590-c660-450f-b35a-841d553fb32d",
                                "6b518e93-b160-47e7-b02d-34d41c9677d3",
                                "635a5aa3-9e41-4e6d-9493-9a49634ecc7a",
                                "f64f4732-d68b-48af-84ce-95cf4c8bb89f",
                                "2ae1c90a-0fc9-403b-8cb0-a9622c51ea67"
                            ],
                            "lastUpdateTime": "2024-05-15T15:29:10.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-05-29T23:27:14.229Z",
                            "objectGuid": "095b6929-44b9-4525-a0cc-9ef4552011f3",
                            "objectSid": "S-1-5-21-3468690955-1439461270-1872542213-1140",
                            "ou": "corp.sgnl.ai/Company/Users",
                            "samAccountName": "Wendolyn.Garber",
                            "servicePrincipalNames": [],
                            "title": null,
                            "upn": "Wendolyn.Garber@sgnldemos.com",
                            "userAccountControl": 512,
                            "userAccountControlFlags": [
                                "NORMAL_ACCOUNT"
                            ]
                        }
                    ],
                    "primaryDisplayName": "Wendolyn Garber"
                },
                {
                    "archived": false,
                    "creationTime": "2024-08-25T18:04:22.000Z",
                    "earliestSeenTraffic": "2024-09-06T14:51:25.118Z",
                    "emailAddresses": [],
                    "entityId": "45dc40e2-7b7b-4f38-9ac7-98f4a35b24e1",
                    "hasADDomainAdminRole": true,
                    "impactScore": 0.98,
                    "inactive": true,
                    "learned": false,
                    "markTime": null,
                    "mostRecentActivity": "2024-09-11T14:07:43.455Z",
                    "riskScore": 0.65,
                    "riskScoreSeverity": "MEDIUM",
                    "riskScoreWithoutLinkedAccounts": 0.6454871339999999,
                    "secondaryDisplayName": "WHOLESALECHIPS.CO\\sgnl-user",
                    "shared": false,
                    "stale": false,
                    "watched": false,
                    "type": "USER",
                    "riskFactors": [
                        {
                            "score": 0.6,
                            "severity": "MEDIUM",
                            "type": "WEAK_PASSWORD_POLICY"
                        },
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "DUPLICATE_PASSWORD"
                        },
                        {
                            "score": 0.15,
                            "severity": "NORMAL",
                            "type": "INACTIVE_ACCOUNT"
                        }
                    ],
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "sgnl-user",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "ecae212e-cad3-4e90-be32-3a2b3fca0cb4",
                                "716f3866-b339-4c25-9085-43460c7f125c",
                                "dde8448e-42a9-4729-bcab-778085c8066d",
                                "dd133c9c-74c5-42af-b446-596f130eee8f",
                                "7d9386af-b875-4264-bc23-2e769b03bc85",
                                "3fe4484c-02f4-43b0-aa51-e1ce10f0ad5c"
                            ],
                            "creationTime": "2024-08-25T18:04:22.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": "Built-in account for administering the computer/domain",
                            "dn": "CN=sgnl-user,CN=Users,DC=wholesalechips,DC=co",
                            "domain": "WHOLESALECHIPS.CO",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "ecae212e-cad3-4e90-be32-3a2b3fca0cb4",
                                "716f3866-b339-4c25-9085-43460c7f125c",
                                "dde8448e-42a9-4729-bcab-778085c8066d",
                                "dd133c9c-74c5-42af-b446-596f130eee8f",
                                "7d9386af-b875-4264-bc23-2e769b03bc85",
                                "3fe4484c-02f4-43b0-aa51-e1ce10f0ad5c",
                                "925b0caa-edbb-46c6-80a0-1700950a7a86",
                                "6d68930f-414e-4f00-85fe-28b868cbb910"
                            ],
                            "lastUpdateTime": "2024-08-25T18:04:22.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-09-11T14:07:43.455Z",
                            "objectGuid": "45dc40e2-7b7b-4f38-9ac7-98f4a35b24e1",
                            "objectSid": "S-1-5-21-1361080754-2191010971-608695987-500",
                            "ou": null,
                            "samAccountName": "sgnl-user",
                            "servicePrincipalNames": [],
                            "title": null,
                            "upn": null,
                            "userAccountControl": 512,
                            "userAccountControlFlags": [
                                "NORMAL_ACCOUNT"
                            ]
                        }
                    ],
                    "primaryDisplayName": "sgnl-user"
                }
            ]
        }
    },
    "extensions": {
        "runTime": 21,
        "remainingPoints": 499998,
        "reset": 9995,
        "consumedPoints": 2
    }
}`

	UserResponsePage2 = `{
    "data": {
        "entities": {
            "pageInfo": {
                "hasNextPage": true,
                "endCursor": "eyJyaXNrU2NvcmUiOjAuNjQwNDc5MTcxNzM1MjQ4OSwiX2lkIjoiODNhNDllZjEtMTdhNy00ZmE0LWI5MGYtOTE0MmRmYTQ5NTc3In0="
            },
            "nodes": [
                {
                    "archived": false,
                    "creationTime": "2024-05-15T15:16:27.000Z",
                    "earliestSeenTraffic": "2024-05-23T02:00:59.187Z",
                    "emailAddresses": [],
                    "entityId": "c1732de2-853c-4375-a479-17b0afbe114f",
                    "hasADDomainAdminRole": true,
                    "impactScore": 0.98,
                    "inactive": false,
                    "learned": true,
                    "markTime": null,
                    "mostRecentActivity": "2024-09-20T22:00:15.650Z",
                    "riskScore": 0.64,
                    "riskScoreSeverity": "MEDIUM",
                    "riskScoreWithoutLinkedAccounts": 0.6427350773921063,
                    "secondaryDisplayName": "CORP.SGNL.AI\\marc",
                    "shared": false,
                    "stale": false,
                    "watched": false,
                    "type": "USER",
                    "riskFactors": [
                        {
                            "score": 0.6,
                            "severity": "MEDIUM",
                            "type": "WEAK_PASSWORD_POLICY"
                        },
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "DUPLICATE_PASSWORD"
                        },
                        {
                            "score": 0.07987954870600919,
                            "severity": "NORMAL",
                            "type": "STALE_ACCOUNT_USAGE"
                        }
                    ],
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "marc",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "5d9de0b3-0fb6-4045-a082-2c091eab7c0c",
                                "e9254f97-959a-4f55-a282-392a49e381d2",
                                "cc1ea590-c660-450f-b35a-841d553fb32d",
                                "6b518e93-b160-47e7-b02d-34d41c9677d3",
                                "635a5aa3-9e41-4e6d-9493-9a49634ecc7a",
                                "5d02ca1a-5201-4f4d-9967-7b44274e8454"
                            ],
                            "creationTime": "2024-05-15T15:16:27.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": "Built-in account for administering the computer/domain",
                            "dn": "CN=marc,CN=Users,DC=corp,DC=sgnl,DC=ai",
                            "domain": "CORP.SGNL.AI",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "5d9de0b3-0fb6-4045-a082-2c091eab7c0c",
                                "e9254f97-959a-4f55-a282-392a49e381d2",
                                "cc1ea590-c660-450f-b35a-841d553fb32d",
                                "6b518e93-b160-47e7-b02d-34d41c9677d3",
                                "635a5aa3-9e41-4e6d-9493-9a49634ecc7a",
                                "5d02ca1a-5201-4f4d-9967-7b44274e8454",
                                "f64f4732-d68b-48af-84ce-95cf4c8bb89f",
                                "2ae1c90a-0fc9-403b-8cb0-a9622c51ea67"
                            ],
                            "lastUpdateTime": "2024-05-15T15:16:27.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-09-20T22:00:15.650Z",
                            "objectGuid": "c1732de2-853c-4375-a479-17b0afbe114f",
                            "objectSid": "S-1-5-21-3468690955-1439461270-1872542213-500",
                            "ou": null,
                            "samAccountName": "marc",
                            "servicePrincipalNames": [],
                            "title": null,
                            "upn": "marc@corp.sgnl.ai",
                            "userAccountControl": 66048,
                            "userAccountControlFlags": [
                                "NORMAL_ACCOUNT",
                                "DONT_EXPIRE_PASSWORD"
                            ]
                        }
                    ],
                    "primaryDisplayName": "marc"
                },
                {
                    "archived": false,
                    "creationTime": "2024-08-25T18:18:00.000Z",
                    "earliestSeenTraffic": "2024-09-04T02:23:23.435Z",
                    "emailAddresses": [],
                    "entityId": "83a49ef1-17a7-4fa4-b90f-9142dfa49577",
                    "hasADDomainAdminRole": true,
                    "impactScore": 0.4,
                    "inactive": false,
                    "learned": false,
                    "markTime": null,
                    "mostRecentActivity": "2024-09-12T15:02:40.094Z",
                    "riskScore": 0.64,
                    "riskScoreSeverity": "MEDIUM",
                    "riskScoreWithoutLinkedAccounts": 0.6404791717713373,
                    "secondaryDisplayName": "WHOLESALECHIPS.CO\\sgnl.sor",
                    "shared": false,
                    "stale": false,
                    "watched": false,
                    "type": "USER",
                    "riskFactors": [
                        {
                            "score": 0.6,
                            "severity": "MEDIUM",
                            "type": "WEAK_PASSWORD_POLICY"
                        },
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "DUPLICATE_PASSWORD"
                        },
                        {
                            "score": 0.02240067243031055,
                            "severity": "NORMAL",
                            "type": "LDAP_RECONNAISSANCE"
                        }
                    ],
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "sgnl sor",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "68bd95ed-9d9f-4ad1-baf3-f2c004b7fd18",
                                "dd133c9c-74c5-42af-b446-596f130eee8f"
                            ],
                            "creationTime": "2024-08-25T18:18:00.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": "Used for SGNL SoR",
                            "dn": "CN=sgnl sor,CN=Users,DC=wholesalechips,DC=co",
                            "domain": "WHOLESALECHIPS.CO",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "68bd95ed-9d9f-4ad1-baf3-f2c004b7fd18",
                                "dd133c9c-74c5-42af-b446-596f130eee8f",
                                "925b0caa-edbb-46c6-80a0-1700950a7a86",
                                "6d68930f-414e-4f00-85fe-28b868cbb910"
                            ],
                            "lastUpdateTime": "2024-08-25T18:18:00.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-09-12T15:02:40.094Z",
                            "objectGuid": "83a49ef1-17a7-4fa4-b90f-9142dfa49577",
                            "objectSid": "S-1-5-21-1361080754-2191010971-608695987-1104",
                            "ou": null,
                            "samAccountName": "sgnl.sor",
                            "servicePrincipalNames": [],
                            "title": null,
                            "upn": "sgnl.sor@wholesalechips.co",
                            "userAccountControl": 66048,
                            "userAccountControlFlags": [
                                "NORMAL_ACCOUNT",
                                "DONT_EXPIRE_PASSWORD"
                            ]
                        }
                    ],
                    "primaryDisplayName": "sgnl sor"
                }
            ]
        }
    },
    "extensions": {
        "runTime": 22,
        "remainingPoints": 499998,
        "reset": 9996,
        "consumedPoints": 2
    }
}`

	UserResponsePage3 = `{
    "data": {
        "entities": {
            "pageInfo": {
                "hasNextPage": false,
                "endCursor": null
            },
            "nodes": [
                {
                    "archived": false,
                    "creationTime": "2024-05-23T15:08:11.000Z",
                    "earliestSeenTraffic": null,
                    "emailAddresses": [],
                    "entityId": "6b4c76ba-2493-4a87-bfb3-1ea91985cce5",
                    "hasADDomainAdminRole": false,
                    "impactScore": 0.25,
                    "inactive": true,
                    "learned": false,
                    "markTime": null,
                    "mostRecentActivity": null,
                    "riskScore": 0.63,
                    "riskScoreSeverity": "MEDIUM",
                    "riskScoreWithoutLinkedAccounts": 0.6262344230994479,
                    "secondaryDisplayName": "CORP.SGNL.AI\\alejandro.bacong",
                    "shared": false,
                    "stale": true,
                    "watched": false,
                    "type": "USER",
                    "riskFactors": [
                        {
                            "score": 0.55,
                            "severity": "MEDIUM",
                            "type": "HAS_ATTACK_PATH"
                        },
                        {
                            "score": 0.4,
                            "severity": "NORMAL",
                            "type": "STALE_ACCOUNT"
                        },
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "WEAK_PASSWORD_POLICY"
                        },
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "DUPLICATE_PASSWORD"
                        }
                    ],
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "Alejandro Bacong",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "cc1ea590-c660-450f-b35a-841d553fb32d"
                            ],
                            "creationTime": "2024-05-23T15:08:11.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": null,
                            "dn": "CN=Alejandro Bacong,OU=Users,OU=Company,DC=corp,DC=sgnl,DC=ai",
                            "domain": "CORP.SGNL.AI",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "cc1ea590-c660-450f-b35a-841d553fb32d",
                                "2ae1c90a-0fc9-403b-8cb0-a9622c51ea67"
                            ],
                            "lastUpdateTime": "2024-05-23T15:08:11.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": null,
                            "objectGuid": "6b4c76ba-2493-4a87-bfb3-1ea91985cce5",
                            "objectSid": "S-1-5-21-3468690955-1439461270-1872542213-2101",
                            "ou": "corp.sgnl.ai/Company/Users",
                            "samAccountName": "alejandro.bacong",
                            "servicePrincipalNames": [],
                            "title": null,
                            "upn": "alejandro.bacong@wholesalechips.co",
                            "userAccountControl": 66048,
                            "userAccountControlFlags": [
                                "NORMAL_ACCOUNT",
                                "DONT_EXPIRE_PASSWORD"
                            ]
                        }
                    ],
                    "primaryDisplayName": "Alejandro Bacong"
                }
            ]
        }
    },
    "extensions": {
        "runTime": 22,
        "remainingPoints": 499999,
        "reset": 9994,
        "consumedPoints": 1
    }
}`

	EndpointResponsePage1 = `{
    "data": {
        "entities": {
            "pageInfo": {
                "hasNextPage": true,
                "endCursor": "eyJyaXNrU2NvcmUiOjAuNDU5NDAwMDAwMDAwMDAwMDMsIl9pZCI6IjNjN2FlYmI5LTQxMWItNGVlOS1iNDgxLWU4ODFmMjlhZmNjOCJ9"
            },
            "nodes": [
                {
                    "agentId": "84a3c4307fee48ef96deeca4a6377cbc",
                    "agentVersion": "7.15.18511.0",
                    "archived": false,
                    "cid": "8693deb4-bf13-4cfb-8855-ee118d9a0243",
                    "creationTime": "2024-05-29T21:30:17.000Z",
                    "earliestSeenTraffic": "2024-05-29T21:33:13.904Z",
                    "entityId": "89be47c3-f51b-48af-884a-ecb02ed0807a",
                    "guestAccountEnabled": false,
                    "hasADDomainAdminRole": false,
                    "hasRole": true,
                    "hostName": "alice-win11.corp.sgnl.ai",
                    "impactScore": 0,
                    "inactive": true,
                    "lastIpAddress": "1.1.1.1",
                    "learned": true,
                    "markTime": null,
                    "mostRecentActivity": "2024-06-18T21:40:54.682Z",
                    "primaryDisplayName": "alice-win11",
                    "riskScore": 0.48,
                    "riskScoreSeverity": "MEDIUM",
                    "secondaryDisplayName": "alice-win11.corp.sgnl.ai",
                    "shared": false,
                    "stale": true,
                    "staticIpAddresses": [],
                    "type": "ENDPOINT",
                    "unmanaged": false,
                    "watched": false,
                    "ztaScore": null,
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "alice-win11",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "b69ca14c-e919-42ba-a21e-62f34c402a13"
                            ],
                            "creationTime": "2024-05-29T21:30:17.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": null,
                            "dn": "CN=alice-win11,OU=Computers,OU=Company,DC=corp,DC=sgnl,DC=ai",
                            "domain": "CORP.SGNL.AI",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "b69ca14c-e919-42ba-a21e-62f34c402a13"
                            ],
                            "lastUpdateTime": "2024-05-29T21:30:17.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-06-18T21:40:54.682Z",
                            "objectGuid": "89be47c3-f51b-48af-884a-ecb02ed0807a",
                            "objectSid": "S-1-5-21-3468690955-1439461270-1872542213-2103",
                            "ou": "corp.sgnl.ai/Company/Computers",
                            "samAccountName": "ALICE-WIN11$",
                            "servicePrincipalNames": [
                                "TERMSRV/ALICE-WIN11",
                                "TERMSRV/alice-win11.corp.sgnl.ai",
                                "RestrictedKrbHost/alice-win11",
                                "HOST/alice-win11",
                                "RestrictedKrbHost/alice-win11.corp.sgnl.ai",
                                "HOST/alice-win11.corp.sgnl.ai"
                            ],
                            "title": null,
                            "upn": null,
                            "userAccountControl": 4096,
                            "userAccountControlFlags": [
                                "WORKSTATION_TRUST_ACCOUNT"
                            ]
                        }
                    ],
                    "riskFactors": [
                        {
                            "score": 0.4,
                            "severity": "NORMAL",
                            "type": "STALE_ACCOUNT"
                        },
                        {
                            "score": 0.4,
                            "severity": "NORMAL",
                            "type": "SMB_SIGNING_DISABLED"
                        }
                    ]
                },
                {
                    "agentId": "eca21da34c934e8e95c97a4f7af1d9a5",
                    "agentVersion": "7.15.18514.0",
                    "archived": false,
                    "cid": "8693deb4-bf13-4cfb-8855-ee118d9a0243",
                    "creationTime": "2024-05-15T15:17:19.000Z",
                    "earliestSeenTraffic": "2024-05-23T02:00:59.187Z",
                    "entityId": "3c7aebb9-411b-4ee9-b481-e881f29afcc8",
                    "guestAccountEnabled": null,
                    "hasADDomainAdminRole": true,
                    "hasRole": true,
                    "hostName": "mj-dc.corp.sgnl.ai",
                    "impactScore": 0,
                    "inactive": false,
                    "lastIpAddress": "1.1.1.1",
                    "learned": true,
                    "markTime": null,
                    "mostRecentActivity": "2024-09-20T22:00:15.650Z",
                    "primaryDisplayName": "mj-dc",
                    "riskScore": 0.46,
                    "riskScoreSeverity": "MEDIUM",
                    "secondaryDisplayName": "mj-dc.corp.sgnl.ai",
                    "shared": false,
                    "stale": false,
                    "staticIpAddresses": [],
                    "type": "ENDPOINT",
                    "unmanaged": false,
                    "watched": true,
                    "ztaScore": 28,
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "mj-dc",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "a3f5d59f-40af-45cd-95ce-19dfdd6c2386",
                                "95cebf5d-36a6-4994-bbdb-693a60e13749",
                                "239dcac1-6d00-4cff-a894-400386750d79"
                            ],
                            "creationTime": "2024-05-15T15:17:19.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": null,
                            "dn": "CN=mj-dc,OU=Domain Controllers,DC=corp,DC=sgnl,DC=ai",
                            "domain": "CORP.SGNL.AI",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "a3f5d59f-40af-45cd-95ce-19dfdd6c2386",
                                "95cebf5d-36a6-4994-bbdb-693a60e13749",
                                "239dcac1-6d00-4cff-a894-400386750d79",
                                "f64f4732-d68b-48af-84ce-95cf4c8bb89f"
                            ],
                            "lastUpdateTime": "2024-05-15T15:17:19.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-09-20T22:00:15.650Z",
                            "objectGuid": "3c7aebb9-411b-4ee9-b481-e881f29afcc8",
                            "objectSid": "S-1-5-21-3468690955-1439461270-1872542213-1000",
                            "ou": "corp.sgnl.ai/Domain Controllers",
                            "samAccountName": "mj-dc$",
                            "servicePrincipalNames": [
                                "Dfsr-12F9A27C-BF97-4787-9364-D31B6C55EB04/mj-dc.corp.sgnl.ai",
                                "ldap/mj-dc.corp.sgnl.ai/ForestDnsZones.corp.sgnl.ai",
                                "ldap/mj-dc.corp.sgnl.ai/DomainDnsZones.corp.sgnl.ai",
                                "TERMSRV/mj-dc",
                                "TERMSRV/mj-dc.corp.sgnl.ai",
                                "DNS/mj-dc.corp.sgnl.ai",
                                "GC/mj-dc.corp.sgnl.ai/corp.sgnl.ai",
                                "RestrictedKrbHost/mj-dc.corp.sgnl.ai",
                                "RestrictedKrbHost/mj-dc",
                                "RPC/a905d6eb-fc70-43e4-b48e-0e4c14822b7e._msdcs.corp.sgnl.ai",
                                "HOST/mj-dc/CORP",
                                "HOST/mj-dc.corp.sgnl.ai/CORP",
                                "HOST/mj-dc",
                                "HOST/mj-dc.corp.sgnl.ai",
                                "HOST/mj-dc.corp.sgnl.ai/corp.sgnl.ai",
                                "E3514235-4B06-11D1-AB04-00C04FC2DCD2/a905d6eb-fc70-43e4-b48e-0e4c14822b7e/corp.sgnl.ai",
                                "ldap/mj-dc/CORP",
                                "ldap/a905d6eb-fc70-43e4-b48e-0e4c14822b7e._msdcs.corp.sgnl.ai",
                                "ldap/mj-dc.corp.sgnl.ai/CORP",
                                "ldap/mj-dc",
                                "ldap/mj-dc.corp.sgnl.ai",
                                "ldap/mj-dc.corp.sgnl.ai/corp.sgnl.ai"
                            ],
                            "title": null,
                            "upn": null,
                            "userAccountControl": 532480,
                            "userAccountControlFlags": [
                                "SERVER_TRUST_ACCOUNT",
                                "TRUSTED_FOR_DELEGATION"
                            ]
                        }
                    ],
                    "riskFactors": [
                        {
                            "score": 0.4,
                            "severity": "NORMAL",
                            "type": "WATCHED"
                        },
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "SPOOLER_SERVICE_RUNNING"
                        }
                    ]
                }
            ]
        }
    },
    "extensions": {
        "runTime": 20,
        "remainingPoints": 499998,
        "reset": 8517,
        "consumedPoints": 2
    }
}`

	EndpointResponsePage2 = `{
    "data": {
        "entities": {
            "pageInfo": {
                "hasNextPage": false,
                "endCursor": "eyJyaXNrU2NvcmUiOjAuMywiX2lkIjoiZmQxZTBmMGItZjFlMS00MjI0LThkNjAtNGYyOTdhYTkxYzI5In0="
            },
            "nodes": [
                {
                    "agentId": "3af65068d68a4440b52bbe1ecacaae14",
                    "agentVersion": "7.15.18514.0",
                    "archived": false,
                    "cid": "8693deb4-bf13-4cfb-8855-ee118d9a0243",
                    "creationTime": "2024-08-25T18:06:23.000Z",
                    "earliestSeenTraffic": "2024-09-04T02:23:23.435Z",
                    "entityId": "fd1e0f0b-f1e1-4224-8d60-4f297aa91c29",
                    "guestAccountEnabled": null,
                    "hasADDomainAdminRole": true,
                    "hasRole": true,
                    "hostName": "se-demo-active-.wholesalechips.co",
                    "impactScore": 0,
                    "inactive": false,
                    "lastIpAddress": "1.1.1.1",
                    "learned": false,
                    "markTime": null,
                    "mostRecentActivity": "2024-09-12T15:02:40.094Z",
                    "primaryDisplayName": "SE-Demo-Active-",
                    "riskScore": 0.3,
                    "riskScoreSeverity": "NORMAL",
                    "secondaryDisplayName": "se-demo-active-.wholesalechips.co",
                    "shared": false,
                    "stale": false,
                    "staticIpAddresses": [],
                    "type": "ENDPOINT",
                    "unmanaged": false,
                    "watched": false,
                    "ztaScore": 28,
                    "accounts": [
                        {
                            "__typename": "ActiveDirectoryAccountDescriptor",
                            "archived": false,
                            "cn": "SE-Demo-Active-",
                            "consistencyGuid": null,
                            "containingGroupIds": [
                                "4c92e0f3-0d13-4de8-8be5-58aed02cd8bd"
                            ],
                            "creationTime": "2024-08-25T18:06:23.000Z",
                            "dataSource": "ACTIVE_DIRECTORY",
                            "department": null,
                            "description": null,
                            "dn": "CN=SE-Demo-Active-,OU=Domain Controllers,DC=wholesalechips,DC=co",
                            "domain": "WHOLESALECHIPS.CO",
                            "enabled": true,
                            "expirationTime": null,
                            "flattenedContainingGroupIds": [
                                "4c92e0f3-0d13-4de8-8be5-58aed02cd8bd",
                                "6d68930f-414e-4f00-85fe-28b868cbb910"
                            ],
                            "lastUpdateTime": "2024-08-25T18:06:23.000Z",
                            "lockoutTime": null,
                            "mostRecentActivity": "2024-09-12T15:02:40.094Z",
                            "objectGuid": "fd1e0f0b-f1e1-4224-8d60-4f297aa91c29",
                            "objectSid": "S-1-5-21-1361080754-2191010971-608695987-1000",
                            "ou": "wholesalechips.co/Domain Controllers",
                            "samAccountName": "SE-Demo-Active-$",
                            "servicePrincipalNames": [
                                "Dfsr-12F9A27C-BF97-4787-9364-D31B6C55EB04/SE-Demo-Active-.wholesalechips.co",
                                "ldap/SE-Demo-Active-.wholesalechips.co/ForestDnsZones.wholesalechips.co",
                                "ldap/SE-Demo-Active-.wholesalechips.co/DomainDnsZones.wholesalechips.co",
                                "TERMSRV/SE-Demo-Active-",
                                "TERMSRV/SE-Demo-Active-.wholesalechips.co",
                                "DNS/SE-Demo-Active-.wholesalechips.co",
                                "GC/SE-Demo-Active-.wholesalechips.co/wholesalechips.co",
                                "RestrictedKrbHost/SE-Demo-Active-.wholesalechips.co",
                                "RestrictedKrbHost/SE-Demo-Active-",
                                "RPC/04b5a0e0-d1c9-43fb-a8a9-e37ddc8100ac._msdcs.wholesalechips.co",
                                "HOST/SE-Demo-Active-/WHOLESALECHIPS",
                                "HOST/SE-Demo-Active-.wholesalechips.co/WHOLESALECHIPS",
                                "HOST/SE-Demo-Active-",
                                "HOST/SE-Demo-Active-.wholesalechips.co",
                                "HOST/SE-Demo-Active-.wholesalechips.co/wholesalechips.co",
                                "E3514235-4B06-11D1-AB04-00C04FC2DCD2/04b5a0e0-d1c9-43fb-a8a9-e37ddc8100ac/wholesalechips.co",
                                "ldap/SE-Demo-Active-/WHOLESALECHIPS",
                                "ldap/04b5a0e0-d1c9-43fb-a8a9-e37ddc8100ac._msdcs.wholesalechips.co",
                                "ldap/SE-Demo-Active-.wholesalechips.co/WHOLESALECHIPS",
                                "ldap/SE-Demo-Active-",
                                "ldap/SE-Demo-Active-.wholesalechips.co",
                                "ldap/SE-Demo-Active-.wholesalechips.co/wholesalechips.co"
                            ],
                            "title": null,
                            "upn": null,
                            "userAccountControl": 532480,
                            "userAccountControlFlags": [
                                "SERVER_TRUST_ACCOUNT",
                                "TRUSTED_FOR_DELEGATION"
                            ]
                        }
                    ],
                    "riskFactors": [
                        {
                            "score": 0.3,
                            "severity": "NORMAL",
                            "type": "SPOOLER_SERVICE_RUNNING"
                        }
                    ]
                }
            ]
        }
    },
    "extensions": {
        "runTime": 20,
        "remainingPoints": 499998,
        "reset": 8337,
        "consumedPoints": 2
    }
}`

	EndpointResponsePage3 = `{
    "data": {
        "entities": {
            "pageInfo": {
                "hasNextPage": false,
                "endCursor": null
            },
            "nodes": []
        }
    },
    "extensions": {
        "runTime": 11,
        "remainingPoints": 499998,
        "reset": 9996,
        "consumedPoints": 2
    }
}`

	IncidentResponsePage1 = `{
    "data": {
        "incidents": {
            "pageInfo": {
                "endCursor": "eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0yMFQwMTo1NToxMC4yNzRaIn0sInNlcXVlbmNlSWQiOjE1fQ==",
                "hasNextPage": true
            },
            "nodes": [
                {
                    "endTime": "2024-09-23T13:00:26.350Z",
                    "incidentId": "INC-16",
                    "lifeCycleStage": "NEW",
                    "markedAsRead": false,
                    "severity": "INFO",
                    "startTime": "2024-09-23T13:00:21.000Z",
                    "type": "UNUSUAL_ACTIVITY",
                    "compromisedEntities": [
                        {
                            "archived": false,
                            "creationTime": "2024-05-15T15:17:19.000Z",
                            "entityId": "3c7aebb9-411b-4ee9-b481-e881f29afcc8",
                            "hasADDomainAdminRole": true,
                            "hasRole": true,
                            "learned": true,
                            "markTime": null,
                            "primaryDisplayName": "mj-dc",
                            "riskScore": 0.46,
                            "riskScoreSeverity": "MEDIUM",
                            "secondaryDisplayName": "mj-dc.corp.sgnl.ai",
                            "type": "ENDPOINT",
                            "watched": true
                        }
                    ],
                    "alertEvents": [
                        {
                            "alertId": "5c941395-4f44-465a-abdd-87b2aececfbe",
                            "alertType": "PrivilegeEscalationAlert",
                            "endTime": "2024-09-23T13:00:59.999Z",
                            "eventId": "jpppv5",
                            "eventLabel": "Privilege escalation (endpoint)",
                            "eventSeverity": "IMPORTANT",
                            "eventType": "ALERT",
                            "patternId": 51131,
                            "resolved": false,
                            "startTime": "2024-09-23T13:00:59.999Z",
                            "timestamp": "2024-09-23T13:00:23.321Z",
                            "entities": [
                                {
                                    "archived": false,
                                    "creationTime": "2024-05-15T15:17:19.000Z",
                                    "entityId": "3c7aebb9-411b-4ee9-b481-e881f29afcc8",
                                    "hasADDomainAdminRole": true,
                                    "hasRole": true,
                                    "learned": true,
                                    "markTime": null,
                                    "primaryDisplayName": "mj-dc",
                                    "riskScore": 0.46,
                                    "riskScoreSeverity": "MEDIUM",
                                    "secondaryDisplayName": "mj-dc.corp.sgnl.ai",
                                    "type": "ENDPOINT",
                                    "watched": true
                                }
                            ]
                        }
                    ]
                },
                {
                    "endTime": "2024-09-20T01:55:10.274Z",
                    "incidentId": "INC-15",
                    "lifeCycleStage": "NEW",
                    "markedAsRead": false,
                    "severity": "INFO",
                    "startTime": "2024-09-20T01:49:27.000Z",
                    "type": "UNUSUAL_ACTIVITY",
                    "compromisedEntities": [
                        {
                            "archived": false,
                            "creationTime": "2024-05-29T20:45:52.000Z",
                            "entityId": "60ee5bb1-805f-46d2-8f3a-9d7cadc52909",
                            "hasADDomainAdminRole": true,
                            "hasRole": true,
                            "learned": true,
                            "markTime": null,
                            "primaryDisplayName": "Alice Wu",
                            "riskScore": 0.61,
                            "riskScoreSeverity": "MEDIUM",
                            "secondaryDisplayName": "CORP.SGNL.AI\\alice",
                            "type": "USER",
                            "watched": false
                        }
                    ],
                    "alertEvents": [
                        {
                            "alertId": "f6816bcd-9e0c-4ea4-8344-03ea6ab58655",
                            "alertType": "StaleAccountUsageAlert",
                            "endTime": "2024-09-20T01:49:59.999Z",
                            "eventId": "jpppvg",
                            "eventLabel": "Use of stale user account",
                            "eventSeverity": "IMPORTANT",
                            "eventType": "ALERT",
                            "patternId": 51130,
                            "resolved": false,
                            "startTime": "2024-09-20T01:49:59.999Z",
                            "timestamp": "2024-09-20T01:50:27.440Z",
                            "entities": [
                                {
                                    "archived": false,
                                    "creationTime": "2024-05-29T20:45:52.000Z",
                                    "entityId": "60ee5bb1-805f-46d2-8f3a-9d7cadc52909",
                                    "hasADDomainAdminRole": true,
                                    "hasRole": true,
                                    "learned": true,
                                    "markTime": null,
                                    "primaryDisplayName": "Alice Wu",
                                    "riskScore": 0.61,
                                    "riskScoreSeverity": "MEDIUM",
                                    "secondaryDisplayName": "CORP.SGNL.AI\\alice",
                                    "type": "USER",
                                    "watched": false
                                }
                            ]
                        },
                        {
                            "alertId": "5f1578fb-505d-448e-9d80-39dca742505b",
                            "alertType": "PrivilegeEscalationAlert",
                            "endTime": "2024-09-20T01:55:59.999Z",
                            "eventId": "jpppv2",
                            "eventLabel": "Privilege escalation (user)",
                            "eventSeverity": "IMPORTANT",
                            "eventType": "ALERT",
                            "patternId": 51113,
                            "resolved": false,
                            "startTime": "2024-09-20T01:55:59.999Z",
                            "timestamp": "2024-09-20T01:55:10.224Z",
                            "entities": [
                                {
                                    "archived": false,
                                    "creationTime": "2024-05-29T20:45:52.000Z",
                                    "entityId": "60ee5bb1-805f-46d2-8f3a-9d7cadc52909",
                                    "hasADDomainAdminRole": true,
                                    "hasRole": true,
                                    "learned": true,
                                    "markTime": null,
                                    "primaryDisplayName": "Alice Wu",
                                    "riskScore": 0.61,
                                    "riskScoreSeverity": "MEDIUM",
                                    "secondaryDisplayName": "CORP.SGNL.AI\\alice",
                                    "type": "USER",
                                    "watched": false
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    },
    "extensions": {
        "runTime": 324
    }
}`

	IncidentResponsePage2 = `{
    "data": {
        "incidents": {
            "pageInfo": {
                "endCursor": "eyJlbmRUaW1lIjp7IiRkYXRlIjoiMjAyNC0wOS0wOVQxNDoyODowNC4wMDhaIn0sInNlcXVlbmNlSWQiOjEzfQ==",
                "hasNextPage": true
            },
            "nodes": [
                {
                    "endTime": "2024-09-20T01:37:17.369Z",
                    "incidentId": "INC-14",
                    "lifeCycleStage": "NEW",
                    "markedAsRead": false,
                    "severity": "INFO",
                    "startTime": "2024-09-20T01:36:36.000Z",
                    "type": "UNUSUAL_ACTIVITY",
                    "compromisedEntities": [
                        {
                            "archived": false,
                            "creationTime": "2024-05-15T15:16:27.000Z",
                            "entityId": "c1732de2-853c-4375-a479-17b0afbe114f",
                            "hasADDomainAdminRole": true,
                            "hasRole": true,
                            "learned": true,
                            "markTime": null,
                            "primaryDisplayName": "marc",
                            "riskScore": 0.64,
                            "riskScoreSeverity": "MEDIUM",
                            "secondaryDisplayName": "CORP.SGNL.AI\\marc",
                            "type": "USER",
                            "watched": false
                        }
                    ],
                    "alertEvents": [
                        {
                            "alertId": "d0d3b02b-1ff5-48cb-9b9a-b2d45e70c26d",
                            "alertType": "StaleAccountUsageAlert",
                            "endTime": "2024-09-20T01:36:59.999Z",
                            "eventId": "jpppva",
                            "eventLabel": "Use of stale user account",
                            "eventSeverity": "IMPORTANT",
                            "eventType": "ALERT",
                            "patternId": 51130,
                            "resolved": false,
                            "startTime": "2024-09-20T01:36:59.999Z",
                            "timestamp": "2024-09-20T01:37:17.245Z",
                            "entities": [
                                {
                                    "archived": false,
                                    "creationTime": "2024-05-15T15:16:27.000Z",
                                    "entityId": "c1732de2-853c-4375-a479-17b0afbe114f",
                                    "hasADDomainAdminRole": true,
                                    "hasRole": true,
                                    "learned": true,
                                    "markTime": null,
                                    "primaryDisplayName": "marc",
                                    "riskScore": 0.64,
                                    "riskScoreSeverity": "MEDIUM",
                                    "secondaryDisplayName": "CORP.SGNL.AI\\marc",
                                    "type": "USER",
                                    "watched": false
                                }
                            ]
                        }
                    ]
                },
                {
                    "endTime": "2024-09-09T14:28:04.008Z",
                    "incidentId": "INC-13",
                    "lifeCycleStage": "NEW",
                    "markedAsRead": false,
                    "severity": "INFO",
                    "startTime": "2024-09-09T14:28:00.000Z",
                    "type": "UNUSUAL_ACTIVITY",
                    "compromisedEntities": [
                        {
                            "archived": false,
                            "creationTime": "2024-05-15T15:17:19.000Z",
                            "entityId": "3c7aebb9-411b-4ee9-b481-e881f29afcc8",
                            "hasADDomainAdminRole": true,
                            "hasRole": true,
                            "learned": true,
                            "markTime": null,
                            "primaryDisplayName": "mj-dc",
                            "riskScore": 0.46,
                            "riskScoreSeverity": "MEDIUM",
                            "secondaryDisplayName": "mj-dc.corp.sgnl.ai",
                            "type": "ENDPOINT",
                            "watched": true
                        }
                    ],
                    "alertEvents": [
                        {
                            "alertId": "6c41cb74-63d6-46b2-a71d-239266737d71",
                            "alertType": "PrivilegeEscalationAlert",
                            "endTime": "2024-09-09T14:28:59.999Z",
                            "eventId": "jpppvf",
                            "eventLabel": "Privilege escalation (endpoint)",
                            "eventSeverity": "IMPORTANT",
                            "eventType": "ALERT",
                            "patternId": 51131,
                            "resolved": false,
                            "startTime": "2024-09-09T14:28:59.999Z",
                            "timestamp": "2024-09-09T14:28:01.520Z",
                            "entities": [
                                {
                                    "archived": false,
                                    "creationTime": "2024-05-15T15:17:19.000Z",
                                    "entityId": "3c7aebb9-411b-4ee9-b481-e881f29afcc8",
                                    "hasADDomainAdminRole": true,
                                    "hasRole": true,
                                    "learned": true,
                                    "markTime": null,
                                    "primaryDisplayName": "mj-dc",
                                    "riskScore": 0.46,
                                    "riskScoreSeverity": "MEDIUM",
                                    "secondaryDisplayName": "mj-dc.corp.sgnl.ai",
                                    "type": "ENDPOINT",
                                    "watched": true
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    },
    "extensions": {
        "runTime": 41
    }
}`

	IncidentResponsePage3 = `{
    "data": {
        "incidents": {
            "pageInfo": {
                "endCursor": null,
                "hasNextPage": false
            },
            "nodes": [
                {
                    "endTime": "2024-09-04T02:30:22.214Z",
                    "incidentId": "INC-12",
                    "lifeCycleStage": "NEW",
                    "markedAsRead": false,
                    "severity": "INFO",
                    "startTime": "2024-09-04T02:23:23.000Z",
                    "type": "SUSPICIOUS_DOMAIN_ACTIVITY",
                    "compromisedEntities": [
                        {
                            "archived": false,
                            "creationTime": "2024-08-25T18:18:00.000Z",
                            "entityId": "83a49ef1-17a7-4fa4-b90f-9142dfa49577",
                            "hasADDomainAdminRole": true,
                            "hasRole": true,
                            "learned": false,
                            "markTime": null,
                            "primaryDisplayName": "sgnl sor",
                            "riskScore": 0.64,
                            "riskScoreSeverity": "MEDIUM",
                            "secondaryDisplayName": "WHOLESALECHIPS.CO\\sgnl.sor",
                            "type": "USER",
                            "watched": false
                        },
                        {
                            "archived": false,
                            "creationTime": "2024-09-04T02:23:23.435Z",
                            "entityId": "40ff0c2d-a1d3-3676-a924-7688b73c163a",
                            "hasADDomainAdminRole": false,
                            "hasRole": false,
                            "learned": false,
                            "markTime": null,
                            "primaryDisplayName": "1.1.1.1",
                            "riskScore": 0.16,
                            "riskScoreSeverity": "NORMAL",
                            "secondaryDisplayName": "",
                            "type": "ENDPOINT",
                            "watched": false
                        }
                    ],
                    "alertEvents": [
                        {
                            "alertId": "0b41c631-6b41-4d8f-abd5-3946aaf45652",
                            "alertType": "LdapReconnaissanceAlert",
                            "endTime": "2024-09-04T02:23:59.999Z",
                            "eventId": "jpppv6",
                            "eventLabel": "Suspicious LDAP search (Kerberos misconfiguration)",
                            "eventSeverity": "IMPORTANT",
                            "eventType": "ALERT",
                            "patternId": 51106,
                            "resolved": false,
                            "startTime": "2024-09-04T02:23:59.999Z",
                            "timestamp": "2024-09-04T02:30:19.695Z",
                            "entities": [
                                {
                                    "archived": false,
                                    "creationTime": "2024-08-25T18:18:00.000Z",
                                    "entityId": "83a49ef1-17a7-4fa4-b90f-9142dfa49577",
                                    "hasADDomainAdminRole": true,
                                    "hasRole": true,
                                    "learned": false,
                                    "markTime": null,
                                    "primaryDisplayName": "sgnl sor",
                                    "riskScore": 0.64,
                                    "riskScoreSeverity": "MEDIUM",
                                    "secondaryDisplayName": "WHOLESALECHIPS.CO\\sgnl.sor",
                                    "type": "USER",
                                    "watched": false
                                },
                                {
                                    "archived": false,
                                    "creationTime": "2024-09-04T02:23:23.435Z",
                                    "entityId": "40ff0c2d-a1d3-3676-a924-7688b73c163a",
                                    "hasADDomainAdminRole": false,
                                    "hasRole": false,
                                    "learned": false,
                                    "markTime": null,
                                    "primaryDisplayName": "1.1.1.1",
                                    "riskScore": 0.16,
                                    "riskScoreSeverity": "NORMAL",
                                    "secondaryDisplayName": "",
                                    "type": "ENDPOINT",
                                    "watched": false
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    },
    "extensions": {
        "runTime": 31
    }
}`

	DetectListResponseFirstPage = `{
        "meta": {
            "query_time": 0.006885612,
            "pagination": {
                "offset": 0,
                "limit": 2,
                "total": 6
            },
            "powered_by": "legacy-detects",
            "trace_id": "07d38d83-5ea4-41d0-9357-b975f516ac54"
        },
        "resources": [
            "ldt:9b9b1e4f7512492f95f8039c065a28a9:4298709414",
            "ldt:9b9b1e4f7512492f95f8039c065a28a9:4298086570"
        ],
        "errors": []
    }`

	DetectListResponseMiddlePage = `{
        "meta": {
            "query_time": 0.00764229,
            "pagination": {
                "offset": 2,
                "limit": 2,
                "total": 6
            },
            "powered_by": "legacy-detects",
            "trace_id": "2415f590-faae-44fb-bf0c-843d5dbb095a"
        },
        "resources": [
            "ldt:9b9b1e4f7512492f95f8039c065a28a9:4295459139",
            "ldt:9b9b1e4f7512492f95f8039c065a28a9:1169567"
        ],
        "errors": [

        ]
    }`

	DetectListResponseLastPage = `{
        "meta": {
            "query_time": 0.007624657,
            "pagination": {
                "offset": 4,
                "limit": 2,
                "total": 6
            },
            "powered_by": "legacy-detects",
            "trace_id": "4ebd37c6-5f1b-4a2c-bcb4-6445d55521cb"
        },
        "resources": [
            "ldt:eca21da34c934e8e95c97a4f7af1d9a5:77310702382",
            "ldt:eca21da34c934e8e95c97a4f7af1d9a5:77309428075"
        ],
        "errors": [

        ]
    }`

	DetectDetailedResponseFirstPage = `{
        "meta": {
            "query_time": 0.008815014,
            "powered_by": "legacy-detects",
            "trace_id": "7288f7be-cd6b-4365-85a8-e05cf7cfd9d8"
        },
        "resources": [
            {
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "created_timestamp": "2025-01-22T19:52:49.295871976Z",
                "detection_id": "ldt:9b9b1e4f7512492f95f8039c065a28a9:4298086570",
                "device": {
                    "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "agent_load_flags": "0",
                    "agent_local_time": "2025-01-22T19:45:19.682Z",
                    "agent_version": "7.19.18913.0",
                    "bios_manufacturer": "Xen",
                    "bios_version": "4.11.amazon",
                    "config_id_base": "65994767",
                    "config_id_build": "18913",
                    "config_id_platform": "3",
                    "external_ip": "1.1.1.1",
                    "hostname": "EC2AMAZ-L4LAU4Q",
                    "first_seen": "2025-01-22T19:37:47Z",
                    "last_login_timestamp": "2025-01-22T19:40:45Z",
                    "last_login_user": "Administrator",
                    "last_seen": "2025-01-22T19:45:30Z",
                    "local_ip": "1.1.1.1",
                    "mac_address": "01-01-01-01-01-01",
                    "machine_domain": "",
                    "major_version": "10",
                    "minor_version": "0",
                    "os_version": "Windows Server 2022",
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "3",
                    "product_type_desc": "Server",
                    "status": "normal",
                    "system_manufacturer": "Xen",
                    "system_product_name": "HVM domU",
                    "modified_timestamp": "2025-01-22T19:50:46Z",
                    "instance_id": "i-04d26bf36004d2941",
                    "service_provider": "AWS_EC2_V2",
                    "service_provider_account_id": "{{OMITTED}}"
                },
                "behaviors": [
                    {
                        "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                        "timestamp": "2025-01-22T19:48:54Z",
                        "template_instance_id": "1343",
                        "behavior_id": "10303",
                        "filename": "cmd.exe",
                        "filepath": "\\Device\\HarddiskVolume1\\Windows\\System32\\cmd.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "cmd.exe  crowdstrike_test_critical",
                        "scenario": "suspicious_activity",
                        "objective": "Falcon Detection Method",
                        "tactic": "Falcon Overwatch",
                        "tactic_id": "CSTA0006",
                        "technique": "Malicious Activity",
                        "technique_id": "CST0002",
                        "display_name": "TestTriggerCritical",
                        "description": "A critical level detection was triggered on this process for testing purposes.",
                        "severity": 90,
                        "confidence": 100,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "Administrator",
                        "user_id": "S-1-5-21-1176167308-4253926863-1726221433-500",
                        "control_graph_id": "ctg:9b9b1e4f7512492f95f8039c065a28a9:4298086570",
                        "triggering_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:4344646281",
                        "sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                        "md5": "448d1a22fb3e4e05dace52091152cc27",
                        "parent_details": {
                            "parent_sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                            "parent_md5": "448d1a22fb3e4e05dace52091152cc27",
                            "parent_cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                            "parent_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:158082492"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    }
                ],
                "email_sent": false,
                "first_behavior": "2025-01-22T19:48:54Z",
                "last_behavior": "2025-01-22T19:48:54Z",
                "max_confidence": 100,
                "max_severity": 90,
                "max_severity_displayname": "Critical",
                "show_in_ui": true,
                "status": "new",
                "hostinfo": {
                    "domain": ""
                },
                "seconds_to_triaged": 0,
                "seconds_to_resolved": 0,
                "behaviors_processed": [
                    "pid:9b9b1e4f7512492f95f8039c065a28a9:4344646281:10303"
                ],
                "date_updated": "2025-01-22T19:53:10Z"
            },
            {
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "created_timestamp": "2025-01-22T21:14:05.663856584Z",
                "detection_id": "ldt:9b9b1e4f7512492f95f8039c065a28a9:4298709414",
                "device": {
                    "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "agent_load_flags": "0",
                    "agent_local_time": "2025-01-22T19:45:19.682Z",
                    "agent_version": "7.19.18913.0",
                    "bios_manufacturer": "Xen",
                    "bios_version": "4.11.amazon",
                    "config_id_base": "65994767",
                    "config_id_build": "18913",
                    "config_id_platform": "3",
                    "external_ip": "1.1.1.1",
                    "hostname": "EC2AMAZ-L4LAU4Q",
                    "first_seen": "2025-01-22T19:37:47Z",
                    "last_login_timestamp": "2025-01-22T19:40:45Z",
                    "last_login_user": "Administrator",
                    "last_seen": "2025-01-22T21:09:28Z",
                    "local_ip": "1.1.1.1",
                    "mac_address": "01-01-01-01-01-01",
                    "machine_domain": "",
                    "major_version": "10",
                    "minor_version": "0",
                    "os_version": "Windows Server 2022",
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "3",
                    "product_type_desc": "Server",
                    "status": "normal",
                    "system_manufacturer": "Xen",
                    "system_product_name": "HVM domU",
                    "modified_timestamp": "2025-01-22T21:13:07Z",
                    "instance_id": "i-04d26bf36004d2941",
                    "service_provider": "AWS_EC2_V2",
                    "service_provider_account_id": "{{OMITTED}}"
                },
                "behaviors": [
                    {
                        "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                        "timestamp": "2025-01-22T21:13:47Z",
                        "template_instance_id": "438",
                        "behavior_id": "82",
                        "filename": "cmd.exe",
                        "filepath": "\\Device\\HarddiskVolume1\\Windows\\System32\\cmd.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                        "scenario": "suspicious_activity",
                        "objective": "Follow Through",
                        "tactic": "Execution",
                        "tactic_id": "TA0002",
                        "technique": "Command and Scripting Interpreter",
                        "technique_id": "T1059",
                        "display_name": "UnexpectedSvchostProcess",
                        "description": "An unexpected process ran svchost.exe. Adversaries can masquerade malware as a system process to evade detection. Review the executable.",
                        "severity": 50,
                        "confidence": 90,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "Administrator",
                        "user_id": "S-1-5-21-1176167308-4253926863-1726221433-500",
                        "control_graph_id": "ctg:9b9b1e4f7512492f95f8039c065a28a9:4298709414",
                        "triggering_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:6086000492",
                        "sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                        "md5": "448d1a22fb3e4e05dace52091152cc27",
                        "parent_details": {
                            "parent_sha256": "26e89cb7afcb09c11b5563c3650196a0f935a95ed44bf1a52c261049db2c4fec",
                            "parent_md5": "",
                            "parent_cmdline": "C:\\Windows\\Explorer.EXE",
                            "parent_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:30879573"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    },
                    {
                        "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                        "timestamp": "2025-01-22T21:13:47Z",
                        "template_instance_id": "438",
                        "behavior_id": "10228",
                        "filename": "powershell.exe",
                        "filepath": "\\Device\\HarddiskVolume1\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "powershell",
                        "scenario": "evade_detection",
                        "objective": "Keep Access",
                        "tactic": "Defense Evasion",
                        "tactic_id": "TA0005",
                        "technique": "Process Hollowing",
                        "technique_id": "T1055.012",
                        "display_name": "SvchostUnexpectedParent",
                        "description": "A service host process launched suspended under an unusual parent. This might indicate a malicious process preparing to inject into svchost for a malicious purpose. Investigate the process tree.",
                        "severity": 70,
                        "confidence": 80,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "Administrator",
                        "user_id": "S-1-5-21-1176167308-4253926863-1726221433-500",
                        "control_graph_id": "ctg:9b9b1e4f7512492f95f8039c065a28a9:4298709414",
                        "triggering_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:6089699227",
                        "sha256": "38f4384643b3fa0de714d2367b712c2e0fa1c89e2cfd131ae6b831ad962b1033",
                        "md5": "dd6f4b7818a253887b8ea86515f6fb7d",
                        "parent_details": {
                            "parent_sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                            "parent_md5": "448d1a22fb3e4e05dace52091152cc27",
                            "parent_cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                            "parent_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:6086000492"
                        },
                        "pattern_disposition": 272,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": true,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": true,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    },
                    {
                        "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                        "timestamp": "2025-01-22T21:13:47Z",
                        "template_instance_id": "438",
                        "behavior_id": "82",
                        "filename": "powershell.exe",
                        "filepath": "\\Device\\HarddiskVolume1\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "powershell",
                        "scenario": "suspicious_activity",
                        "objective": "Follow Through",
                        "tactic": "Execution",
                        "tactic_id": "TA0002",
                        "technique": "Command and Scripting Interpreter",
                        "technique_id": "T1059",
                        "display_name": "UnexpectedSvchostProcess",
                        "description": "An unexpected process ran svchost.exe. Adversaries can masquerade malware as a system process to evade detection. Review the executable.",
                        "severity": 50,
                        "confidence": 90,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "Administrator",
                        "user_id": "S-1-5-21-1176167308-4253926863-1726221433-500",
                        "control_graph_id": "ctg:9b9b1e4f7512492f95f8039c065a28a9:4298709414",
                        "triggering_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:6089699227",
                        "sha256": "38f4384643b3fa0de714d2367b712c2e0fa1c89e2cfd131ae6b831ad962b1033",
                        "md5": "dd6f4b7818a253887b8ea86515f6fb7d",
                        "parent_details": {
                            "parent_sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                            "parent_md5": "448d1a22fb3e4e05dace52091152cc27",
                            "parent_cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                            "parent_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:6086000492"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    }
                ],
                "email_sent": false,
                "first_behavior": "2025-01-22T21:13:47Z",
                "last_behavior": "2025-01-22T21:13:47Z",
                "max_confidence": 90,
                "max_severity": 70,
                "max_severity_displayname": "High",
                "show_in_ui": true,
                "status": "new",
                "hostinfo": {
                    "domain": ""
                },
                "seconds_to_triaged": 0,
                "seconds_to_resolved": 0,
                "behaviors_processed": [
                    "pid:9b9b1e4f7512492f95f8039c065a28a9:6089699227:10228",
                    "pid:9b9b1e4f7512492f95f8039c065a28a9:6089699227:82",
                    "pid:9b9b1e4f7512492f95f8039c065a28a9:6086000492:82"
                ],
                "date_updated": "2025-01-22T21:14:27Z"
            }
        ],
        "errors": [

        ]
    }`

	DetectDetailedResponseMiddlePage = `{
        "meta": {
            "query_time": 0.007343097,
            "powered_by": "legacy-detects",
            "trace_id": "c5d13812-743a-42bb-b169-7bd486af2975"
        },
        "resources": [
            {
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "created_timestamp": "2025-01-22T19:48:16.921558696Z",
                "detection_id": "ldt:9b9b1e4f7512492f95f8039c065a28a9:1169567",
                "device": {
                    "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "agent_load_flags": "0",
                    "agent_local_time": "2025-01-22T19:45:19.682Z",
                    "agent_version": "7.19.18913.0",
                    "bios_manufacturer": "Xen",
                    "bios_version": "4.11.amazon",
                    "config_id_base": "65994767",
                    "config_id_build": "18913",
                    "config_id_platform": "3",
                    "external_ip": "1.1.1.1",
                    "hostname": "EC2AMAZ-L4LAU4Q",
                    "first_seen": "2025-01-22T19:37:47Z",
                    "last_login_timestamp": "2025-01-22T19:40:45Z",
                    "last_login_user": "Administrator",
                    "last_seen": "2025-01-22T19:45:30Z",
                    "local_ip": "1.1.1.1",
                    "mac_address": "01-01-01-01-01-01",
                    "machine_domain": "",
                    "major_version": "10",
                    "minor_version": "0",
                    "os_version": "Windows Server 2022",
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "3",
                    "product_type_desc": "Server",
                    "status": "normal",
                    "system_manufacturer": "Xen",
                    "system_product_name": "HVM domU",
                    "modified_timestamp": "2025-01-22T19:47:13Z",
                    "instance_id": "i-04d26bf36004d2941",
                    "service_provider": "AWS_EC2_V2",
                    "service_provider_account_id": "{{OMITTED}}"
                },
                "behaviors": [
                    {
                        "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                        "timestamp": "2025-01-22T19:44:18Z",
                        "template_instance_id": "1342",
                        "behavior_id": "10304",
                        "filename": "cmd.exe",
                        "filepath": "\\Device\\HarddiskVolume1\\Windows\\System32\\cmd.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "cmd.exe  crowdstrike_test_high",
                        "scenario": "suspicious_activity",
                        "objective": "Falcon Detection Method",
                        "tactic": "Falcon Overwatch",
                        "tactic_id": "CSTA0006",
                        "technique": "Malicious Activity",
                        "technique_id": "CST0002",
                        "display_name": "TestTriggerHigh",
                        "description": "A high level detection was triggered on this process for testing purposes.",
                        "severity": 70,
                        "confidence": 100,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "Administrator",
                        "user_id": "S-1-5-21-1176167308-4253926863-1726221433-500",
                        "control_graph_id": "ctg:9b9b1e4f7512492f95f8039c065a28a9:1169567",
                        "triggering_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:166798888",
                        "sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                        "md5": "448d1a22fb3e4e05dace52091152cc27",
                        "parent_details": {
                            "parent_sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                            "parent_md5": "448d1a22fb3e4e05dace52091152cc27",
                            "parent_cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                            "parent_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:158082492"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    }
                ],
                "email_sent": false,
                "first_behavior": "2025-01-22T19:44:18Z",
                "last_behavior": "2025-01-22T19:44:18Z",
                "max_confidence": 100,
                "max_severity": 70,
                "max_severity_displayname": "High",
                "show_in_ui": true,
                "status": "new",
                "hostinfo": {
                    "domain": ""
                },
                "seconds_to_triaged": 0,
                "seconds_to_resolved": 0,
                "behaviors_processed": [
                    "pid:9b9b1e4f7512492f95f8039c065a28a9:166798888:10304"
                ],
                "date_updated": "2025-01-22T19:48:38Z"
            },
            {
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "created_timestamp": "2025-01-22T19:48:57.812628852Z",
                "detection_id": "ldt:9b9b1e4f7512492f95f8039c065a28a9:4295459139",
                "device": {
                    "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "agent_load_flags": "0",
                    "agent_local_time": "2025-01-22T19:45:19.682Z",
                    "agent_version": "7.19.18913.0",
                    "bios_manufacturer": "Xen",
                    "bios_version": "4.11.amazon",
                    "config_id_base": "65994767",
                    "config_id_build": "18913",
                    "config_id_platform": "3",
                    "external_ip": "1.1.1.1",
                    "hostname": "EC2AMAZ-L4LAU4Q",
                    "first_seen": "2025-01-22T19:37:47Z",
                    "last_login_timestamp": "2025-01-22T19:40:45Z",
                    "last_login_user": "Administrator",
                    "last_seen": "2025-01-22T19:45:30Z",
                    "local_ip": "1.1.1.1",
                    "mac_address": "01-01-01-01-01-01",
                    "machine_domain": "",
                    "major_version": "10",
                    "minor_version": "0",
                    "os_version": "Windows Server 2022",
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "3",
                    "product_type_desc": "Server",
                    "status": "normal",
                    "system_manufacturer": "Xen",
                    "system_product_name": "HVM domU",
                    "modified_timestamp": "2025-01-22T19:47:13Z",
                    "instance_id": "i-04d26bf36004d2941",
                    "service_provider": "AWS_EC2_V2",
                    "service_provider_account_id": "{{OMITTED}}"
                },
                "behaviors": [
                    {
                        "device_id": "9b9b1e4f7512492f95f8039c065a28a9",
                        "timestamp": "2025-01-22T19:48:48Z",
                        "template_instance_id": "1343",
                        "behavior_id": "10303",
                        "filename": "cmd.exe",
                        "filepath": "\\Device\\HarddiskVolume1\\Windows\\System32\\cmd.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "cmd.exe  crowdstrike_test_critical",
                        "scenario": "suspicious_activity",
                        "objective": "Falcon Detection Method",
                        "tactic": "Falcon Overwatch",
                        "tactic_id": "CSTA0006",
                        "technique": "Malicious Activity",
                        "technique_id": "CST0002",
                        "display_name": "TestTriggerCritical",
                        "description": "A critical level detection was triggered on this process for testing purposes.",
                        "severity": 90,
                        "confidence": 100,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "Administrator",
                        "user_id": "S-1-5-21-1176167308-4253926863-1726221433-500",
                        "control_graph_id": "ctg:9b9b1e4f7512492f95f8039c065a28a9:4295459139",
                        "triggering_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:4341293422",
                        "sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                        "md5": "448d1a22fb3e4e05dace52091152cc27",
                        "parent_details": {
                            "parent_sha256": "41871dade953d9f40f4aa445fc19982ab59d263c8aa93d7f67a1451663a09a57",
                            "parent_md5": "448d1a22fb3e4e05dace52091152cc27",
                            "parent_cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                            "parent_process_graph_id": "pid:9b9b1e4f7512492f95f8039c065a28a9:158082492"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    }
                ],
                "email_sent": false,
                "first_behavior": "2025-01-22T19:48:48Z",
                "last_behavior": "2025-01-22T19:48:48Z",
                "max_confidence": 100,
                "max_severity": 90,
                "max_severity_displayname": "Critical",
                "show_in_ui": true,
                "status": "new",
                "hostinfo": {
                    "domain": ""
                },
                "seconds_to_triaged": 0,
                "seconds_to_resolved": 0,
                "behaviors_processed": [
                    "pid:9b9b1e4f7512492f95f8039c065a28a9:4341293422:10303"
                ],
                "date_updated": "2025-01-22T19:49:19Z"
            }
        ],
        "errors": [

        ]
    }`

	DetectDetailedResponseLastPage = `{
        "meta": {
            "query_time": 0.011262357,
            "powered_by": "legacy-detects",
            "trace_id": "192a1a20-2784-4ea8-8cfc-34a211f2902a"
        },
        "resources": [
            {
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "created_timestamp": "2024-12-05T02:25:10.199790849Z",
                "detection_id": "ldt:eca21da34c934e8e95c97a4f7af1d9a5:77310702382",
                "device": {
                    "device_id": "eca21da34c934e8e95c97a4f7af1d9a5",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "agent_load_flags": "0",
                    "agent_local_time": "2024-12-05T02:19:45.022Z",
                    "agent_version": "7.17.18721.0",
                    "bios_manufacturer": "Microsoft Corporation",
                    "bios_version": "Hyper-V UEFI Release v4.1",
                    "config_id_base": "65994763",
                    "config_id_build": "18721",
                    "config_id_platform": "3",
                    "external_ip": "1.1.1.1",
                    "hostname": "mj-dc",
                    "first_seen": "2024-12-05T02:16:13Z",
                    "last_login_timestamp": "2024-12-05T02:24:14Z",
                    "last_login_user": "jan.f",
                    "last_seen": "2024-12-05T02:20:06Z",
                    "local_ip": "1.1.1.1",
                    "mac_address": "01-01-01-01-01-01",
                    "machine_domain": "corp.sgnl.ai",
                    "major_version": "10",
                    "minor_version": "0",
                    "os_version": "Windows Server 2022",
                    "ou": [
                        "Domain Controllers"
                    ],
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "2",
                    "product_type_desc": "Domain Controller",
                    "site_name": "Default-First-Site-Name",
                    "status": "normal",
                    "system_manufacturer": "Microsoft Corporation",
                    "system_product_name": "Virtual Machine",
                    "groups": [
                        "2a8b900d486e4e9eaa024723d6f3742a"
                    ],
                    "modified_timestamp": "2024-12-05T02:24:15Z",
                    "instance_id": "4220508a-d2a1-466f-9187-40594db3256b",
                    "service_provider": "AZURE",
                    "service_provider_account_id": "{{OMITTED}}"
                },
                "behaviors": [
                    {
                        "device_id": "eca21da34c934e8e95c97a4f7af1d9a5",
                        "timestamp": "2024-12-05T02:25:00Z",
                        "template_instance_id": "438",
                        "behavior_id": "10228",
                        "filename": "powershell.exe",
                        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\" ",
                        "scenario": "evade_detection",
                        "objective": "Keep Access",
                        "tactic": "Defense Evasion",
                        "tactic_id": "TA0005",
                        "technique": "Process Hollowing",
                        "technique_id": "T1055.012",
                        "display_name": "SvchostUnexpectedParent",
                        "description": "A service host process launched suspended under an unusual parent. This might indicate a malicious process preparing to inject into svchost for a malicious purpose. Investigate the process tree.",
                        "severity": 70,
                        "confidence": 80,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "jan.f",
                        "user_id": "S-1-5-21-3468690955-1439461270-1872542213-7102",
                        "control_graph_id": "ctg:eca21da34c934e8e95c97a4f7af1d9a5:77310702382",
                        "triggering_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481146653568",
                        "sha256": "75d6634a676fb0bea5bfd8d424e2bd4f685f3885853637ea143b2671a3bb76e9",
                        "md5": "0bc8a4cd1e07390bafd741e1fc0399a3",
                        "parent_details": {
                            "parent_sha256": "26b1a027ba0581ae6448c03a4c842f6d94b672f4c3024aabd8993c64bc181163",
                            "parent_md5": "4ed94002301ee4ae46ddf33e076c8dba",
                            "parent_cmdline": "\"C:\\Windows\\system32\\RunDll32.exe\" C:\\Windows\\System32\\SHELL32.dll,RunAsNewUser_RunDLL Local\\{4ddb9f3f-700c-4bd6-9fc0-eaf85c01d25b}.00002568",
                            "parent_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481144985587"
                        },
                        "pattern_disposition": 272,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": true,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": true,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    },
                    {
                        "device_id": "eca21da34c934e8e95c97a4f7af1d9a5",
                        "timestamp": "2024-12-05T02:25:00Z",
                        "template_instance_id": "438",
                        "behavior_id": "82",
                        "filename": "powershell.exe",
                        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\" ",
                        "scenario": "suspicious_activity",
                        "objective": "Follow Through",
                        "tactic": "Execution",
                        "tactic_id": "TA0002",
                        "technique": "Command and Scripting Interpreter",
                        "technique_id": "T1059",
                        "display_name": "UnexpectedSvchostProcess",
                        "description": "An unexpected process ran svchost.exe. Adversaries can masquerade malware as a system process to evade detection. Review the executable.",
                        "severity": 50,
                        "confidence": 90,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "jan.f",
                        "user_id": "S-1-5-21-3468690955-1439461270-1872542213-7102",
                        "control_graph_id": "ctg:eca21da34c934e8e95c97a4f7af1d9a5:77310702382",
                        "triggering_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481146653568",
                        "sha256": "75d6634a676fb0bea5bfd8d424e2bd4f685f3885853637ea143b2671a3bb76e9",
                        "md5": "0bc8a4cd1e07390bafd741e1fc0399a3",
                        "parent_details": {
                            "parent_sha256": "26b1a027ba0581ae6448c03a4c842f6d94b672f4c3024aabd8993c64bc181163",
                            "parent_md5": "4ed94002301ee4ae46ddf33e076c8dba",
                            "parent_cmdline": "\"C:\\Windows\\system32\\RunDll32.exe\" C:\\Windows\\System32\\SHELL32.dll,RunAsNewUser_RunDLL Local\\{4ddb9f3f-700c-4bd6-9fc0-eaf85c01d25b}.00002568",
                            "parent_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481144985587"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    }
                ],
                "email_sent": false,
                "first_behavior": "2024-12-05T02:25:00Z",
                "last_behavior": "2024-12-05T02:25:00Z",
                "max_confidence": 90,
                "max_severity": 70,
                "max_severity_displayname": "High",
                "show_in_ui": true,
                "status": "new",
                "hostinfo": {
                    "active_directory_dn_display": [
                        "Domain Controllers"
                    ],
                    "domain": ""
                },
                "seconds_to_triaged": 0,
                "seconds_to_resolved": 0,
                "behaviors_processed": [
                    "pid:eca21da34c934e8e95c97a4f7af1d9a5:481146653568:82",
                    "pid:eca21da34c934e8e95c97a4f7af1d9a5:481146653568:10228"
                ],
                "date_updated": "2024-12-05T02:25:16Z"
            },
            {
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "created_timestamp": "2024-12-05T02:25:09.838903415Z",
                "detection_id": "ldt:eca21da34c934e8e95c97a4f7af1d9a5:77309428075",
                "device": {
                    "device_id": "eca21da34c934e8e95c97a4f7af1d9a5",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "agent_load_flags": "0",
                    "agent_local_time": "2024-12-05T02:19:45.022Z",
                    "agent_version": "7.17.18721.0",
                    "bios_manufacturer": "Microsoft Corporation",
                    "bios_version": "Hyper-V UEFI Release v4.1",
                    "config_id_base": "65994763",
                    "config_id_build": "18721",
                    "config_id_platform": "3",
                    "external_ip": "1.1.1.1",
                    "hostname": "mj-dc",
                    "first_seen": "2024-12-05T02:16:13Z",
                    "last_login_timestamp": "2024-12-05T02:24:14Z",
                    "last_login_user": "jan.f",
                    "last_seen": "2024-12-05T02:20:06Z",
                    "local_ip": "1.1.1.1",
                    "mac_address": "01-01-01-01-01-01",
                    "machine_domain": "corp.sgnl.ai",
                    "major_version": "10",
                    "minor_version": "0",
                    "os_version": "Windows Server 2022",
                    "ou": [
                        "Domain Controllers"
                    ],
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "2",
                    "product_type_desc": "Domain Controller",
                    "site_name": "Default-First-Site-Name",
                    "status": "normal",
                    "system_manufacturer": "Microsoft Corporation",
                    "system_product_name": "Virtual Machine",
                    "groups": [
                        "2a8b900d486e4e9eaa024723d6f3742a"
                    ],
                    "modified_timestamp": "2024-12-05T02:24:15Z",
                    "instance_id": "4220508a-d2a1-466f-9187-40594db3256b",
                    "service_provider": "AZURE",
                    "service_provider_account_id": "{{OMITTED}}"
                },
                "behaviors": [
                    {
                        "device_id": "eca21da34c934e8e95c97a4f7af1d9a5",
                        "timestamp": "2024-12-05T02:24:31Z",
                        "template_instance_id": "1342",
                        "behavior_id": "10304",
                        "filename": "cmd.exe",
                        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "\"C:\\Windows\\system32\\cmd.exe\" crowdstrike_test_high",
                        "scenario": "suspicious_activity",
                        "objective": "Falcon Detection Method",
                        "tactic": "Falcon Overwatch",
                        "tactic_id": "CSTA0006",
                        "technique": "Malicious Activity",
                        "technique_id": "CST0002",
                        "display_name": "TestTriggerHigh",
                        "description": "A high level detection was triggered on this process for testing purposes.",
                        "severity": 70,
                        "confidence": 100,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "jan.f",
                        "user_id": "S-1-5-21-3468690955-1439461270-1872542213-7102",
                        "control_graph_id": "ctg:eca21da34c934e8e95c97a4f7af1d9a5:77309428075",
                        "triggering_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481151226282",
                        "sha256": "54724f38ff2f85c3ff91de434895668b6f39008fc205a668ab6aafad6fb4d93d",
                        "md5": "503ee109ce5cac4bd61084cb28fbd200",
                        "parent_details": {
                            "parent_sha256": "75d6634a676fb0bea5bfd8d424e2bd4f685f3885853637ea143b2671a3bb76e9",
                            "parent_md5": "0bc8a4cd1e07390bafd741e1fc0399a3",
                            "parent_cmdline": "\"C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe\" ",
                            "parent_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481146653568"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    },
                    {
                        "device_id": "eca21da34c934e8e95c97a4f7af1d9a5",
                        "timestamp": "2024-12-05T02:24:38Z",
                        "template_instance_id": "1342",
                        "behavior_id": "10304",
                        "filename": "cmd.exe",
                        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                        "alleged_filetype": "exe",
                        "cmdline": "cmd  crowdstrike_test_high",
                        "scenario": "suspicious_activity",
                        "objective": "Falcon Detection Method",
                        "tactic": "Falcon Overwatch",
                        "tactic_id": "CSTA0006",
                        "technique": "Malicious Activity",
                        "technique_id": "CST0002",
                        "display_name": "TestTriggerHigh",
                        "description": "A high level detection was triggered on this process for testing purposes.",
                        "severity": 70,
                        "confidence": 100,
                        "ioc_type": "",
                        "ioc_value": "",
                        "ioc_source": "",
                        "ioc_description": "",
                        "user_name": "jan.f",
                        "user_id": "S-1-5-21-3468690955-1439461270-1872542213-7102",
                        "control_graph_id": "ctg:eca21da34c934e8e95c97a4f7af1d9a5:77309428075",
                        "triggering_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481157102881",
                        "sha256": "54724f38ff2f85c3ff91de434895668b6f39008fc205a668ab6aafad6fb4d93d",
                        "md5": "503ee109ce5cac4bd61084cb28fbd200",
                        "parent_details": {
                            "parent_sha256": "54724f38ff2f85c3ff91de434895668b6f39008fc205a668ab6aafad6fb4d93d",
                            "parent_md5": "503ee109ce5cac4bd61084cb28fbd200",
                            "parent_cmdline": "\"C:\\Windows\\system32\\cmd.exe\" crowdstrike_test_high",
                            "parent_process_graph_id": "pid:eca21da34c934e8e95c97a4f7af1d9a5:481151226282"
                        },
                        "pattern_disposition": 0,
                        "pattern_disposition_details": {
                            "indicator": false,
                            "detect": false,
                            "inddet_mask": false,
                            "sensor_only": false,
                            "rooting": false,
                            "kill_process": false,
                            "kill_subprocess": false,
                            "quarantine_machine": false,
                            "quarantine_file": false,
                            "policy_disabled": false,
                            "kill_parent": false,
                            "operation_blocked": false,
                            "process_blocked": false,
                            "registry_operation_blocked": false,
                            "critical_process_disabled": false,
                            "bootup_safeguard_enabled": false,
                            "fs_operation_blocked": false,
                            "handle_operation_downgraded": false,
                            "kill_action_failed": false,
                            "blocking_unsupported_or_disabled": false,
                            "suspend_process": false,
                            "suspend_parent": false
                        }
                    }
                ],
                "email_sent": false,
                "first_behavior": "2024-12-05T02:24:31Z",
                "last_behavior": "2024-12-05T02:24:38Z",
                "max_confidence": 100,
                "max_severity": 70,
                "max_severity_displayname": "High",
                "show_in_ui": true,
                "status": "new",
                "hostinfo": {
                    "active_directory_dn_display": [
                        "Domain Controllers"
                    ],
                    "domain": ""
                },
                "seconds_to_triaged": 0,
                "seconds_to_resolved": 0,
                "behaviors_processed": [
                    "pid:eca21da34c934e8e95c97a4f7af1d9a5:481157102881:10304",
                    "pid:eca21da34c934e8e95c97a4f7af1d9a5:481151226282:10304"
                ],
                "date_updated": "2024-12-05T02:25:31Z"
            }
        ],
        "errors": [

        ]
    }`

	DetectResponseSpecializedErr = `{
        "meta": {
            "query_time": 1.64e-7,
            "powered_by": "crowdstrike-api-gateway",
            "trace_id": "968dc340-a865-4065-a4d3-e6ecd94dea74"
        },
        "errors": [
            {
                "code": 404,
                "message": "404: Page Not Found"
            }
        ]
    }`

	AlertListResponseFirstPage = `{
        "meta": {
            "query_time": 0.006885612,
            "pagination": {
                "offset": 0,
                "limit": 2,
                "total": 6
            },
            "powered_by": "alerts",
            "trace_id": "07d38d83-5ea4-41d0-9357-b975f516ac54"
        },
        "resources": [
            "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7086C",
            "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7087D"
        ],
        "errors": []
    }`

	AlertListResponseMiddlePage = `{
        "meta": {
            "query_time": 0.00764229,
            "pagination": {
                "offset": 2,
                "limit": 2,
                "total": 6
            },
            "powered_by": "alerts",
            "trace_id": "2415f590-faae-44fb-bf0c-843d5dbb095a"
        },
        "resources": [
            "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7088E",
            "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7089F"
        ],
        "errors": []
    }`

	AlertListResponseLastPage = `{
        "meta": {
            "query_time": 0.007013842,
            "pagination": {
                "offset": 4,
                "limit": 2,
                "total": 6
            },
            "powered_by": "alerts",
            "trace_id": "8d91fcdc-5e9b-448a-b4ba-8ece5a647ba7"
        },
        "resources": [
            "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7090A",
            "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7091B"
        ],
        "errors": []
    }`

	AlertDetailedResponseFirstPage = `{
        "meta": {
            "query_time": 0.015098675,
            "writes": {
                "resources_affected": 0
            },
            "powered_by": "detectsapi",
            "trace_id": "98c92a82-25eb-4fb8-9055-100d10df2d6e"
        },
        "resources": [
            {
                "added_privileges": [
                    "AdministratorsRole"
                ],
                "aggregate_id": "aggind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7086C",
                "cid": "1234567890abcdef1234567890abcdef",
                "composite_id": "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7086C",
                "confidence": 20,
                "context_timestamp": "2025-08-01T19:22:34.233Z",
                "crawled_timestamp": "2025-08-01T20:22:36.344689563Z",
                "created_timestamp": "2025-08-01T19:23:36.547858841Z",
                "data_domains": [
                    "Identity"
                ],
                "description": "A user received new privileges",
                "display_name": "Privilege escalation (user)",
                "end_time": "2025-08-01T19:22:34.233Z",
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/identity-protection/detections/1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7086C?_cid=test123456789abcdef0123456789abc",
                "id": "ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7086C",
                "mitre_attack": [
                    {
                        "pattern_id": 51113,
                        "tactic_id": "TA0004",
                        "technique_id": "T1078",
                        "tactic": "Privilege Escalation",
                        "technique": "Valid Accounts"
                    }
                ],
                "name": "IdpEntityPrivilegeEscalationUser",
                "objective": "Gain Access",
                "pattern_id": 51113,
                "poly_id": "AACGk960vxNM-4hV7hGNmgJDLHuOENREmtsZGwxsXx_rpwAATiHWsiR7jY1eZfzHO9v5bTBiDJpIT3MtfdjRLhQS-VD_0g==",
                "previous_privileges": "0",
                "privileges": "2177",
                "product": "idp",
                "scenario": "privilege_escalation",
                "seconds_to_resolved": 0,
                "seconds_to_triaged": 0,
                "severity": 10,
                "severity_name": "Informational",
                "show_in_ui": true,
                "source_account_domain": "ACMECORP.AI",
                "source_account_name": "alejandro.bacong",
                "source_account_object_guid": "9160F713-360C-4331-AC61-EBBAD741A7BD",
                "source_account_object_sid": "S-1-5-21-2931850618-1476300705-2742956860-1105",
                "source_account_sam_account_name": "alejandro.bacong",
                "source_account_upn": "alejandro.bacong@acmecorp.ai",
                "source_products": [
                    "Falcon Identity Protection"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "start_time": "2025-08-01T19:22:34.233Z",
                "status": "new",
                "tactic": "Privilege Escalation",
                "tactic_id": "TA0004",
                "technique": "Valid Accounts",
                "technique_id": "T1078",
                "timestamp": "2025-08-01T19:22:34.848Z",
                "type": "idp-user-endpoint-app-info",
                "updated_timestamp": "2025-08-01T20:22:36.344680541Z",
                "user_name": "alejandro.bacong"
            },
            {
                "aggregate_id": "aggind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7087D",
                "cid": "1234567890abcdef1234567890abcdef",
                "composite_id": "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7087D",
                "confidence": 30,
                "context_timestamp": "2025-08-01T18:45:12.156Z",
                "crawled_timestamp": "2025-08-01T19:45:14.234567890Z",
                "created_timestamp": "2025-08-01T18:46:14.123456789Z",
                "data_domains": [
                    "Endpoint"
                ],
                "description": "Suspicious file execution detected",
                "display_name": "Malicious behavior (endpoint)",
                "end_time": "2025-08-01T18:45:12.156Z",
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/identity-protection/detections/1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7087D?_cid=test123456789abcdef0123456789abc",
                "id": "ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7087D",
                "mitre_attack": [
                    {
                        "pattern_id": 50234,
                        "tactic_id": "TA0002",
                        "technique_id": "T1204",
                        "tactic": "Execution",
                        "technique": "User Execution"
                    }
                ],
                "name": "MaliciousFileExecution",
                "objective": "Execute Code",
                "pattern_id": 50234,
                "product": "edr",
                "scenario": "malicious_execution",
                "seconds_to_resolved": 0,
                "seconds_to_triaged": 0,
                "severity": 6,
                "severity_name": "Medium",
                "show_in_ui": true,
                "source_account_domain": "TESTDOMAIN.COM",
                "source_account_name": "john.doe",
                "source_account_object_guid": "A260F812-470D-5442-BD72-FCCAD851B8CE",
                "source_products": [
                    "Falcon Endpoint Protection"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "start_time": "2025-08-01T18:45:12.156Z",
                "status": "new",
                "tactic": "Execution",
                "tactic_id": "TA0002",
                "technique": "User Execution",
                "technique_id": "T1204",
                "timestamp": "2025-08-01T18:45:12.678Z",
                "type": "edr-endpoint-detection",
                "updated_timestamp": "2025-08-01T19:45:14.234567123Z",
                "user_name": "john.doe"
            }
        ]
    }`

	AlertDetailedResponseMiddlePage = `{
        "meta": {
            "query_time": 0.025680677,
            "writes": {
                "resources_affected": 0
            },
            "powered_by": "detectsapi",
            "trace_id": "2415f590-faae-44fb-bf0c-843d5dbb095a"
        },
        "resources": [
            {
                "aggregate_id": "aggind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7088E",
                "cid": "1234567890abcdef1234567890abcdef",
                "composite_id": "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7088E",
                "confidence": 85,
                "context_timestamp": "2025-08-02T10:15:22.445Z",
                "crawled_timestamp": "2025-08-02T11:15:24.567890123Z",
                "created_timestamp": "2025-08-02T10:16:24.891234567Z",
                "data_domains": [
                    "Network"
                ],
                "description": "Suspicious network communication detected",
                "display_name": "Command and Control (network)",
                "end_time": "2025-08-02T10:15:22.445Z",
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/identity-protection/detections/1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7088E?_cid=test123456789abcdef0123456789abc",
                "id": "ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7088E",
                "mitre_attack": [
                    {
                        "pattern_id": 52445,
                        "tactic_id": "TA0011",
                        "technique_id": "T1071",
                        "tactic": "Command and Control",
                        "technique": "Application Layer Protocol"
                    }
                ],
                "name": "SuspiciousNetworkCommunication",
                "objective": "Command and Control",
                "pattern_id": 52445,
                "product": "edr",
                "scenario": "command_and_control",
                "seconds_to_resolved": 0,
                "seconds_to_triaged": 0,
                "severity": 8,
                "severity_name": "High",
                "show_in_ui": true,
                "source_account_domain": "ENTERPRISE.LOCAL",
                "source_account_name": "service.account",
                "source_account_object_guid": "C370F923-581E-6553-CE83-GDDED962C9DF",
                "source_products": [
                    "Falcon Endpoint Protection"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "start_time": "2025-08-02T10:15:22.445Z",
                "status": "in_progress",
                "tactic": "Command and Control",
                "tactic_id": "TA0011",
                "technique": "Application Layer Protocol",
                "technique_id": "T1071",
                "timestamp": "2025-08-02T10:15:22.889Z",
                "type": "edr-network-detection",
                "updated_timestamp": "2025-08-02T11:15:24.567890456Z",
                "user_name": "service.account"
            },
            {
                "aggregate_id": "aggind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7089F",
                "cid": "1234567890abcdef1234567890abcdef",
                "composite_id": "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7089F",
                "confidence": 90,
                "context_timestamp": "2025-08-02T12:30:45.678Z",
                "crawled_timestamp": "2025-08-02T13:30:47.890123456Z",
                "created_timestamp": "2025-08-02T12:31:47.234567890Z",
                "data_domains": [
                    "Endpoint"
                ],
                "description": "Critical malware execution detected",
                "display_name": "Malware execution (endpoint)",
                "end_time": "2025-08-02T12:30:45.678Z",
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/identity-protection/detections/1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7089F?_cid=test123456789abcdef0123456789abc",
                "id": "ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7089F",
                "mitre_attack": [
                    {
                        "pattern_id": 53667,
                        "tactic_id": "TA0005",
                        "technique_id": "T1055",
                        "tactic": "Defense Evasion",
                        "technique": "Process Injection"
                    }
                ],
                "name": "MalwareExecution",
                "objective": "Execute Malware",
                "pattern_id": 53667,
                "product": "edr",
                "scenario": "malware_execution",
                "seconds_to_resolved": 0,
                "seconds_to_triaged": 120,
                "severity": 9,
                "severity_name": "Critical",
                "show_in_ui": true,
                "source_account_domain": "WORKGROUP",
                "source_account_name": "admin.user",
                "source_account_object_guid": "D481FA34-692F-7664-DF94-HGFEF073DAEG",
                "source_products": [
                    "Falcon Endpoint Protection"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "start_time": "2025-08-02T12:30:45.678Z",
                "status": "closed",
                "tactic": "Defense Evasion",
                "tactic_id": "TA0005",
                "technique": "Process Injection",
                "technique_id": "T1055",
                "timestamp": "2025-08-02T12:30:46.123Z",
                "type": "edr-malware-detection",
                "updated_timestamp": "2025-08-02T13:30:47.890123789Z",
                "user_name": "admin.user"
            }
        ]
    }`

	AlertDetailedResponseLastPage = `{
        "meta": {
            "query_time": 0.018754329,
            "writes": {
                "resources_affected": 0
            },
            "powered_by": "detectsapi",
            "trace_id": "8d91fcdc-5e9b-448a-b4ba-8ece5a647ba7"
        },
        "resources": [
            {
                "aggregate_id": "aggind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7090A",
                "cid": "1234567890abcdef1234567890abcdef",
                "composite_id": "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7090A",
                "confidence": 75,
                "context_timestamp": "2025-08-03T14:20:33.567Z",
                "crawled_timestamp": "2025-08-03T15:20:35.789012345Z",
                "created_timestamp": "2025-08-03T14:21:35.345678901Z",
                "data_domains": [
                    "Identity"
                ],
                "description": "User account compromised",
                "display_name": "Account compromise (identity)",
                "end_time": "2025-08-03T14:20:33.567Z",
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/identity-protection/detections/1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7090A?_cid=test123456789abcdef0123456789abc",
                "id": "ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7090A",
                "mitre_attack": [
                    {
                        "pattern_id": 54778,
                        "tactic_id": "TA0006",
                        "technique_id": "T1110",
                        "tactic": "Credential Access",
                        "technique": "Brute Force"
                    }
                ],
                "name": "AccountCompromise",
                "objective": "Credential Access",
                "pattern_id": 54778,
                "product": "idp",
                "scenario": "credential_access",
                "seconds_to_resolved": 300,
                "seconds_to_triaged": 60,
                "severity": 7,
                "severity_name": "High",
                "show_in_ui": true,
                "source_account_domain": "CORPORATE.NET",
                "source_account_name": "security.analyst",
                "source_account_object_guid": "E592GB45-703G-8775-EGA5-IHGFG184EBFH",
                "source_products": [
                    "Falcon Identity Protection"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "start_time": "2025-08-03T14:20:33.567Z",
                "status": "resolved",
                "tactic": "Credential Access",
                "tactic_id": "TA0006",
                "technique": "Brute Force",
                "technique_id": "T1110",
                "timestamp": "2025-08-03T14:20:34.012Z",
                "type": "idp-account-compromise",
                "updated_timestamp": "2025-08-03T15:20:35.789012678Z",
                "user_name": "security.analyst"
            },
            {
                "aggregate_id": "aggind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7091B",
                "cid": "1234567890abcdef1234567890abcdef",
                "composite_id": "1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7091B",
                "confidence": 65,
                "context_timestamp": "2025-08-03T16:45:28.123Z",
                "crawled_timestamp": "2025-08-03T17:45:30.456789012Z",
                "created_timestamp": "2025-08-03T16:46:30.789012345Z",
                "data_domains": [
                    "Endpoint"
                ],
                "description": "Suspicious registry modification",
                "display_name": "Registry tampering (endpoint)",
                "end_time": "2025-08-03T16:45:28.123Z",
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/identity-protection/detections/1234567890abcdef1234567890abcdef:ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7091B?_cid=test123456789abcdef0123456789abc",
                "id": "ind:1234567890abcdef1234567890abcdef:B75E9689-C82E-4EC9-B972-E807CFE7091B",
                "mitre_attack": [
                    {
                        "pattern_id": 55889,
                        "tactic_id": "TA0005",
                        "technique_id": "T1112",
                        "tactic": "Defense Evasion",
                        "technique": "Modify Registry"
                    }
                ],
                "name": "RegistryTampering",
                "objective": "Defense Evasion",
                "pattern_id": 55889,
                "product": "edr",
                "scenario": "registry_modification",
                "seconds_to_resolved": 0,
                "seconds_to_triaged": 0,
                "severity": 4,
                "severity_name": "Medium",
                "show_in_ui": true,
                "source_account_domain": "TESTLAB.ORG",
                "source_account_name": "test.user",
                "source_account_object_guid": "F6A3HC56-814H-9886-FHB6-JIHGH295FCGI",
                "source_products": [
                    "Falcon Endpoint Protection"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "start_time": "2025-08-03T16:45:28.123Z",
                "status": "new",
                "tactic": "Defense Evasion",
                "tactic_id": "TA0005",
                "technique": "Modify Registry",
                "technique_id": "T1112",
                "timestamp": "2025-08-03T16:45:28.567Z",
                "type": "edr-registry-detection",
                "updated_timestamp": "2025-08-03T17:45:30.456789345Z",
                "user_name": "test.user"
            }
        ]
    }`

	// Combined Alerts API responses.
	CombinedAlertResponseFirstPage = `{
		"meta": {
			"query_time": 0.043633983,
			"pagination": {
				"total": 23,
				"limit": 2,
				"after": "eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NjExMTU3MjIxLCJ0ZXN0aWQ6aW5kOjUzODhjNTkyMTg5NDQ0YWQ5ZTg0ZGYwNzFjOGYzOTU0Ojk3ODI3ODI2MTQtMTAzMDMtMzE4MzE1NjgiXSwidG90YWxfZmV0Y2hlZCI6Mn0="
			},
			"powered_by": "detectsapi",
			"trace_id": "00774ad1-9fbf-43c2-9d9e-d90ea9216ba0"
		},
		"errors": [],
		"resources": [
			{
				"agent_id": "test1234567890123456789012345678",
				"aggregate_id": "aggind:test1234567890123456789012345678:625985642613668398",
				"cid": "testcid1234567890123456789012345",
				"cloud_indicator": "false",
				"cmdline": "cat",
				"composite_id": "testcid1234567890123456789012345:ind:test1234567890123456789012345678:625985642593750673-20151-7049",
				"confidence": 80,
				"context_timestamp": "2025-06-16T18:46:54.925Z",
				"control_graph_id": "ctg:test1234567890123456789012345678:625985642613668398",
				"crawled_timestamp": "2025-06-16T19:46:56.218776698Z",
				"created_timestamp": "2025-06-16T18:47:56.231503572Z",
				"data_domains": ["Endpoint"],
				"description": "A process has written a known EICAR test file. Review the files written by the triggered process.",
				"device": {
					"agent_load_flags": "0",
					"agent_local_time": "2025-01-05T13:11:22.705Z",
					"agent_version": "7.24.19504.0",
					"cid": "testcid1234567890123456789012345",
					"config_id_base": "65994763",
					"config_id_build": "19504",
					"config_id_platform": "4",
					"device_id": "test1234567890123456789012345678",
					"external_ip": "192.168.1.100",
					"first_seen": "2025-06-16T17:30:34Z",
					"hostinfo": {"domain": ""},
					"hostname": "TEST-HOST.localdomain",
					"last_seen": "2025-06-16T19:36:52Z",
					"local_ip": "192.168.1.15",
					"mac_address": "00-11-22-33-44-55",
					"machine_domain": "",
					"major_version": "24",
					"minor_version": "5",
					"modified_timestamp": "2025-06-16T19:44:11Z",
					"os_version": "Sequoia (15)",
					"ou": null,
					"platform_id": "1",
					"platform_name": "Mac",
					"product_type_desc": "Workstation",
					"status": "normal",
					"system_manufacturer": "Apple Inc.",
					"system_product_name": "MacBookPro18,1"
				},
				"display_name": "EICARTestFileWrittenMac",
				"email_sent": true,
				"falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/testcid1234567890123456789012345:ind:test1234567890123456789012345678:625985642593750673-20151-7049?_cid=testcid",
				"filename": "cat",
				"filepath": "/bin/cat",
				"global_prevalence": "common",
				"id": "ind:test1234567890123456789012345678:625985642593750673-20151-7049",
				"local_prevalence": "unique",
				"local_process_id": "4077",
				"md5": "27549b2faf16aed0ea0a74bd35c6694e",
				"name": "EICARTestFileWrittenMac",
				"objective": "Follow Through",
				"pattern_id": 20151,
				"platform": "Mac",
				"priority_value": 10,
				"process_id": "625985642593750673",
				"product": "epp",
				"severity": 10,
				"severity_name": "Informational",
				"sha1": "0000000000000000000000000000000000000000",
				"sha256": "9844bd3abe6c59d9cc1ddba3ba35aab555d5ec483db326c183af3c4a6caefbc1",
				"show_in_ui": true,
				"status": "new",
				"tactic": "Execution",
				"tactic_id": "TA0002",
				"technique": "User Execution",
				"technique_id": "T1204",
				"timestamp": "2025-06-16T18:46:54.988Z",
				"type": "ldt",
				"updated_timestamp": "2025-06-16T19:46:56.218762664Z",
				"user_id": "S-1-5-21-1234567890-1234567890-1234567890-2004",
				"user_name": "testuser",
				"user_principal": "testuser@TEST-HOST.localdomain"
			},
			{
				"agent_id": "test2345678901234567890123456789",
				"aggregate_id": "aggind:test2345678901234567890123456789:726196753724579846",
				"cid": "testcid1234567890123456789012345",
				"composite_id": "testcid1234567890123456789012345:ind:test2345678901234567890123456789:726196753704661621-20152-7050",
				"confidence": 75,
				"created_timestamp": "2025-06-17T14:32:18.445672891Z",
				"data_domains": ["Endpoint"],
				"description": "Suspicious network activity detected from unknown source.",
				"display_name": "SuspiciousNetworkActivity",
				"id": "ind:test2345678901234567890123456789:726196753704661621-20152-7050",
				"platform": "Windows",
				"severity": 50,
				"severity_name": "Medium",
				"status": "in_progress",
				"type": "ldt"
			}
		]
	}`

	CombinedAlertResponseMiddlePage = `{
		"meta": {
			"query_time": 0.028765432,
			"pagination": {
				"total": 23,
				"limit": 2,
				"after": "eyJ2ZXJzaW9uIjoidjEiLCJ0b3RhbF9oaXRzIjoyMywidG90YWxfcmVsYXRpb24iOiJlcSIsImNsdXN0ZXJfaWQiOiJ0ZXN0IiwiYWZ0ZXIiOlsxNzQ5NTEyMzQ1Njc4LCJ0ZXN0aWQ6aW5kOmU0NTY3ODkwMTIzNDU2YWI3ODkwY2RlZjEyMzQ1Njc4LTIwMTUzLTcwNTEiXSwidG90YWxfZmV0Y2hlZCI6NH0="
			},
			"powered_by": "detectsapi",
			"trace_id": "11885be2-0gcd-54d3-ae0f-e91fb0327cb1"
		},
		"errors": [],
		"resources": [
			{
				"agent_id": "test3456789012345678901234567890",
				"aggregate_id": "aggind:test3456789012345678901234567890:827307864835690957",
				"cid": "testcid1234567890123456789012345",
				"composite_id": "testcid1234567890123456789012345:ind:test3456789012345678901234567890:827307864815772732-20153-7052",
				"confidence": 90,
				"created_timestamp": "2025-06-18T09:15:42.678912345Z",
				"data_domains": ["Endpoint"],
				"description": "Advanced persistent threat activity successfully mitigated.",
				"display_name": "APTActivity",
				"falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/testcid1234567890123456789012345:ind:test3456789012345678901234567890:827307864815772732-20153-7052?_cid=testcid",
				"id": "ind:test3456789012345678901234567890:827307864815772732-20153-7052",
				"platform": "Linux",
				"priority_value": 80,
				"severity": 80,
				"severity_name": "High",
				"show_in_ui": true,
				"status": "closed",
				"tactic": "Defense Evasion",
				"tactic_id": "TA0005",
				"technique": "Masquerading",
				"technique_id": "T1036",
				"timestamp": "2025-06-18T09:15:42.678Z",
				"type": "ldt",
				"updated_timestamp": "2025-06-18T09:45:12.123456789Z"
			}
		]
	}`

	CombinedAlertResponseLastPage = `{
		"meta": {
			"query_time": 0.015432109,
			"pagination": {
				"total": 23,
				"limit": 2,
				"after": ""
			},
			"powered_by": "detectsapi",
			"trace_id": "22996cf3-1hde-65e4-bf1g-f02gc1438dc2"
		},
		"errors": [],
		"resources": [
			{
				"agent_id": "test4567890123456789012345678901",
				"aggregate_id": "aggind:test4567890123456789012345678901:928418975946802068",
				"cid": "testcid1234567890123456789012345",
				"composite_id": "testcid1234567890123456789012345:ind:test4567890123456789012345678901:928418975926883843-20154-7053",
				"confidence": 95,
				"created_timestamp": "2025-06-19T16:42:18.789123456Z",
				"data_domains": ["Endpoint"],
				"description": "Ransomware indicators detected and quarantined successfully.",
				"display_name": "RansomwareDetection",
				"falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/testcid1234567890123456789012345:ind:test4567890123456789012345678901:928418975926883843-20154-7053?_cid=testcid",
				"id": "ind:test4567890123456789012345678901:928418975926883843-20154-7053",
				"platform": "Windows",
				"priority_value": 100,
				"severity": 100,
				"severity_name": "Critical",
				"show_in_ui": true,
				"status": "new",
				"tactic": "Impact",
				"tactic_id": "TA0040",
				"technique": "Data Encrypted for Impact",
				"technique_id": "T1486",
				"timestamp": "2025-06-19T16:42:18.789Z",
				"type": "ldt",
				"updated_timestamp": "2025-06-19T17:01:22.456789123Z"
			}
		]
	}`
)
