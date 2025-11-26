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

	// Alerts API responses.
	AlertResponseFirstPage = `{
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
      "agent_id": "c36c42b64ce54b39a32e1d57240704c8",
      "aggregate_id": "aggind:c36c42b64ce54b39a32e1d57240704c8:625985642613668398",
      "cid": "8693deb4bf134cfb8855ee118d9a0243",
      "cloud_indicator": "false",
      "cmdline": "cat",
      "composite_id": "8693deb4bf134cfb8855ee118d9a0243:ind:c36c42b64ce54b39a32e1d57240704c8:625985642593750673-20151-7049",
      "confidence": 80,
      "context_timestamp": "2025-06-16T18:46:54.925Z",
      "control_graph_id": "ctg:c36c42b64ce54b39a32e1d57240704c8:625985642613668398",
      "crawled_timestamp": "2025-06-16T19:46:56.218776698Z",
      "created_timestamp": "2025-06-16T18:47:56.231503572Z",
      "data_domains": ["Endpoint"],
      "description": "A process has written a known EICAR test file. Review the files written by the triggered process.",
      "device": {
        "agent_load_flags": "0",
        "agent_local_time": "2610-01-05T13:11:22.705Z",
        "agent_version": "7.24.19504.0",
        "cid": "8693deb4bf134cfb8855ee118d9a0243",
        "config_id_base": "65994763",
        "config_id_build": "19504",
        "config_id_platform": "4",
        "device_id": "c36c42b64ce54b39a32e1d57240704c8",
        "external_ip": "209.122.93.17",
        "first_seen": "2025-06-16T17:30:34Z",
        "hostinfo": {
          "domain": ""
        },
        "hostname": "MLX40LWGRK.localdomain",
        "last_seen": "2025-06-16T19:36:52Z",
        "local_ip": "192.168.1.15",
        "mac_address": "4e-93-8a-29-49-64",
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
      "falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/8693deb4bf134cfb8855ee118d9a0243:ind:c36c42b64ce54b39a32e1d57240704c8:625985642593750673-20151-7049?_cid=g04000s6h3lrs7encahh67joerwbbsje",
      "filename": "cat",
      "filepath": "/bin/cat",
      "files_accessed": [
        {
          "filename": "cat",
          "filepath": "/bin/",
          "timestamp": "1750099615"
        },
        {
          "filename": "zshnW4W3l",
          "filepath": "/private/tmp/",
          "timestamp": "1750099615"
        }
      ],
      "files_written": [
        {
          "filename": "eicar.com",
          "filepath": "/Users/joe/Desktop/",
          "timestamp": "1750099614"
        },
        {
          "filename": "zshnW4W3l",
          "filepath": "/private/tmp/",
          "timestamp": "1750099614"
        }
      ],
      "global_prevalence": "common",
      "id": "ind:c36c42b64ce54b39a32e1d57240704c8:625985642593750673-20151-7049",
      "incident": {
        "created": "2025-06-16T18:46:54Z",
        "end": "2025-06-16T18:46:54Z",
        "id": "inc:c36c42b64ce54b39a32e1d57240704c8:2343cc087cce4683a3c96b49a1e8c863",
        "score": "8.59047945539694",
        "start": "2025-06-16T18:26:55Z"
      },
      "indicator_id": "ind:c36c42b64ce54b39a32e1d57240704c8:625985642593750673-20151-7049",
      "ioc_context": [],
      "local_prevalence": "unique",
      "local_process_id": "4077",
      "md5": "27549b2faf16aed0ea0a74bd35c6694e",
      "mitre_attack": [
        {
          "pattern_id": 20151,
          "tactic_id": "TA0002",
          "technique_id": "T1204",
          "tactic": "Execution",
          "technique": "User Execution"
        }
      ],
      "name": "EICARTestFileWrittenMac",
      "objective": "Follow Through",
      "parent_details": {
        "cmdline": "/bin/zsh -l",
        "filename": "zsh",
        "filepath": "/bin/zsh",
        "local_process_id": "2887",
        "md5": "822dfa404c17133a55fa86c8160b0a1c",
        "process_graph_id": "pid:c36c42b64ce54b39a32e1d57240704c8:625985535262558233",
        "process_id": "625985535262558233",
        "sha256": "7da41b0c36d724529a0769f48635459bbb66cd021f12d16e11245afeb52c4937",
        "timestamp": "2025-06-16T18:47:35Z",
        "user_graph_id": "uid:c36c42b64ce54b39a32e1d57240704c8:502",
        "user_id": "S-1-5-21-4013840944-2679614602-1372758803-2004",
        "user_name": "joe"
      },
      "parent_process_id": "625985535262558233",
      "pattern_disposition": 0,
      "pattern_disposition_description": "Detection, standard detection.",
      "pattern_disposition_details": {
        "blocking_unsupported_or_disabled": false,
        "bootup_safeguard_enabled": false,
        "containment_file_system": false,
        "critical_process_disabled": false,
        "detect": false,
        "fs_operation_blocked": false,
        "handle_operation_downgraded": false,
        "inddet_mask": false,
        "indicator": false,
        "kill_action_failed": false,
        "kill_parent": false,
        "kill_process": false,
        "kill_subprocess": false,
        "mfa_required": false,
        "operation_blocked": false,
        "policy_disabled": false,
        "prevention_provisioning_enabled": false,
        "process_blocked": false,
        "quarantine_file": false,
        "quarantine_machine": false,
        "registry_operation_blocked": false,
        "response_action_already_applied": false,
        "response_action_failed": false,
        "response_action_triggered": false,
        "rooting": false,
        "sensor_only": false,
        "suspend_parent": false,
        "suspend_process": false
      },
      "pattern_id": 20151,
      "platform": "Mac",
      "poly_id": "AACGk960vxNM-4hV7hGNmgJDOQuL6SeY1glLi1xxStXoNQAATiEmNBGGHdYjbo4Hiun25IaeWlNeHeyGfDLO3FYg9vTTuw==",
      "priority_explanation": [
        "[MOD] The severity of the detection: Informational"
      ],
      "priority_value": 10,
      "process_end_time": "1750099614",
      "process_id": "625985642593750673",
      "process_start_time": "1750099614",
      "product": "epp",
      "scenario": "attacker_methodology",
      "seconds_to_resolved": 0,
      "seconds_to_triaged": 0,
      "severity": 10,
      "severity_name": "Informational",
      "sha1": "0000000000000000000000000000000000000000",
      "sha256": "9844bd3abe6c59d9cc1ddba3ba35aab555d5ec483db326c183af3c4a6caefbc1",
      "show_in_ui": true,
      "source_products": ["Falcon Insight"],
      "source_vendors": ["CrowdStrike"],
      "status": "new",
      "tactic": "Execution",
      "tactic_id": "TA0002",
      "technique": "User Execution",
      "technique_id": "T1204",
      "template_instance_id": "7990",
      "timestamp": "2025-06-16T18:46:54.988Z",
      "tree_id": "625985642613668398",
      "tree_root": "625985642593750673",
      "triggering_process_graph_id": "pid:c36c42b64ce54b39a32e1d57240704c8:625985642593750673",
      "type": "ldt",
      "updated_timestamp": "2025-06-16T19:46:56.218762664Z",
      "user_id": "S-1-5-21-4013840944-2679614602-1372758803-2004",
      "user_name": "joe",
      "user_principal": "joe@MLX40LWGRK.localdomain"
    },
    {
      "agent_id": "5388c592189444ad9e84df071c8f3954",
      "aggregate_id": "aggind:5388c592189444ad9e84df071c8f3954:8592364792",
      "alleged_filetype": "exe",
      "cid": "8693deb4bf134cfb8855ee118d9a0243",
      "cloud_indicator": "false",
      "cmdline": "cmd  crowdstrike_test_high",
      "composite_id": "8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:12119912898-10304-117513744",
      "confidence": 100,
      "context_timestamp": "2025-06-13T02:01:54.602Z",
      "control_graph_id": "ctg:5388c592189444ad9e84df071c8f3954:8592364792",
      "crawled_timestamp": "2025-06-13T03:01:57.180370024Z",
      "created_timestamp": "2025-06-13T02:02:57.209949732Z",
      "data_domains": ["Endpoint"],
      "description": "A high level detection was triggered on this process for testing purposes.",
      "device": {
        "agent_load_flags": "1",
        "agent_local_time": "2025-06-10T14:00:03.712Z",
        "agent_version": "7.24.19607.0",
        "bios_manufacturer": "Microsoft Corporation",
        "bios_version": "Hyper-V UEFI Release v4.1",
        "cid": "8693deb4bf134cfb8855ee118d9a0243",
        "config_id_base": "65994767",
        "config_id_build": "19607",
        "config_id_platform": "3",
        "device_id": "5388c592189444ad9e84df071c8f3954",
        "external_ip": "20.83.184.209",
        "first_seen": "2025-06-10T01:49:09Z",
        "groups": ["2a8b900d486e4e9eaa024723d6f3742a"],
        "hostinfo": {
          "active_directory_dn_display": ["Domain Controllers"],
          "domain": "normcorp.ai"
        },
        "hostname": "SGNL-CRWD-Proto",
        "instance_id": "3bfaa67d-0dbd-4d49-8a0a-1cb2b0d2e1af",
        "last_seen": "2025-06-13T03:00:17Z",
        "local_ip": "10.3.0.4",
        "mac_address": "00-0d-3a-56-11-c1",
        "machine_domain": "normcorp.ai",
        "major_version": "10",
        "minor_version": "0",
        "modified_timestamp": "2025-06-13T03:01:35Z",
        "os_version": "Windows Server 2022",
        "ou": ["Domain Controllers"],
        "platform_id": "0",
        "platform_name": "Windows",
        "product_type": "2",
        "product_type_desc": "Domain Controller",
        "service_provider": "AZURE",
        "service_provider_account_id": "9405b466-dc55-4b34-a424-2059ff303a68",
        "site_name": "Default-First-Site-Name",
        "status": "normal",
        "system_manufacturer": "Microsoft Corporation",
        "system_product_name": "Virtual Machine"
      },
      "display_name": "TestTriggerHigh",
      "email_sent": true,
      "falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:12119912898-10304-117513744?_cid=g04000s6h3lrs7encahh67joerwbbsje",
      "filename": "cmd.exe",
      "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
      "global_prevalence": "common",
      "files_accessed": [
        {
          "filename": "cat",
          "filepath": "/bin/",
          "timestamp": "1750099615"
        },
        {
          "filename": "zshnW4W3l",
          "filepath": "/private/tmp/",
          "timestamp": "1750099615"
        }
      ],
      "files_written": [
        {
          "filename": "eicar.com",
          "filepath": "/Users/joe/Desktop/",
          "timestamp": "1750099614"
        },
        {
          "filename": "zshnW4W3l",
          "filepath": "/private/tmp/",
          "timestamp": "1750099614"
        }
      ],
      "grandparent_details": {
        "cmdline": "cmd  crowdstrike_test_high",
        "filename": "cmd.exe",
        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
        "local_process_id": "8868",
        "md5": "22cdd5b627bed769783a7928efed69ae",
        "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10557972040",
        "process_id": "10557972040",
        "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
        "timestamp": "2025-06-11T17:13:16Z",
        "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_name": "user6"
      },
      "incident": {
        "created": "2025-06-16T18:46:54Z",
        "end": "2025-06-16T18:46:54Z",
        "id": "inc:c36c42b64ce54b39a32e1d57240704c8:2343cc087cce4683a3c96b49a1e8c864",
        "score": "8.59047945539694",
        "start": "2025-06-16T18:26:55Z"
      },
      "id": "ind:5388c592189444ad9e84df071c8f3954:12119912898-10304-117513744",
      "indicator_id": "ind:5388c592189444ad9e84df071c8f3954:12119912898-10304-117513744",
      "ioc_context": [],
      "local_prevalence": "low",
      "local_process_id": "3044",
      "logon_domain": "normcorp",
      "md5": "22cdd5b627bed769783a7928efed69ae",
      "mitre_attack": [
        {
          "pattern_id": 10304,
          "tactic_id": "CSTA0006",
          "technique_id": "CST0002",
          "tactic": "Falcon Overwatch",
          "technique": "Malicious Activity"
        }
      ],
      "name": "DemoHighPattern",
      "objective": "Falcon Detection Method",
      "parent_details": {
        "cmdline": "cmd  crowdstrike_test_high",
        "filename": "cmd.exe",
        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
        "local_process_id": "8000",
        "md5": "22cdd5b627bed769783a7928efed69ae",
        "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10653769300",
        "process_id": "10653769300",
        "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
        "timestamp": "2025-06-11T18:39:26Z",
        "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_name": "user6"
      },
      "parent_process_id": "10653769300",
      "pattern_disposition": 0,
      "pattern_disposition_description": "Detection, standard detection.",
      "pattern_disposition_details": {
        "blocking_unsupported_or_disabled": false,
        "bootup_safeguard_enabled": false,
        "containment_file_system": false,
        "critical_process_disabled": false,
        "detect": false,
        "fs_operation_blocked": false,
        "handle_operation_downgraded": false,
        "inddet_mask": false,
        "indicator": false,
        "kill_action_failed": false,
        "kill_parent": false,
        "kill_process": false,
        "kill_subprocess": false,
        "mfa_required": false,
        "operation_blocked": false,
        "policy_disabled": false,
        "prevention_provisioning_enabled": false,
        "process_blocked": false,
        "quarantine_file": false,
        "quarantine_machine": false,
        "registry_operation_blocked": false,
        "response_action_already_applied": false,
        "response_action_failed": false,
        "response_action_triggered": false,
        "rooting": false,
        "sensor_only": false,
        "suspend_parent": false,
        "suspend_process": false
      },
      "pattern_id": 10304,
      "platform": "Windows",
      "poly_id": "AACGk960vxNM-4hV7hGNmgJDGMt_f4XsjSvL9TV99OnITAAATiFkg-CA2Yr0NE5EuxICl1TEUCcCRNBbBp_d2KLZk9JttA==",
      "priority_explanation": [
        "[MOD] The detection is based on Pattern 10304: A high level detection was triggered on this process for testing purposes.",
        "[MOD] The disposition for the detection: no preventative action was taken",
        "[MOD] The parent process was identified as: cmd.exe"
      ],
      "priority_value": 10,
      "process_id": "12119912898",
      "process_start_time": "1749780114",
      "product": "epp",
      "resolution": "ignored",
      "scenario": "suspicious_activity",
      "seconds_to_resolved": 0,
      "seconds_to_triaged": 0,
      "severity": 70,
      "severity_name": "High",
      "sha1": "0000000000000000000000000000000000000000",
      "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
      "show_in_ui": true,
      "source_products": ["Falcon Insight"],
      "source_vendors": ["CrowdStrike"],
      "status": "closed",
      "tactic": "Falcon Overwatch",
      "tactic_id": "CSTA0006",
      "tags": ["ignored"],
      "technique": "Malicious Activity",
      "technique_id": "CST0002",
      "template_instance_id": "1342",
      "timestamp": "2025-06-13T02:01:55.147Z",
      "tree_id": "8592364792",
      "tree_root": "10557972040",
      "triggering_process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:12119912898",
      "type": "ldt",
      "updated_timestamp": "2025-06-18T20:23:35.340932619Z",
      "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
      "user_name": "user6"
    }
  ]
}
`

	AlertResponseMiddlePage = `{
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
      "agent_id": "5388c592189444ad9e84df071c8f3954",
      "aggregate_id": "aggind:5388c592189444ad9e84df071c8f3954:8592364792",
      "alleged_filetype": "exe",
      "cid": "8693deb4bf134cfb8855ee118d9a0243",
      "cloud_indicator": "false",
      "cmdline": "cmd  crowdstrike_test_high",
      "composite_id": "8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10653769300-10304-81908752",
      "confidence": 100,
      "context_timestamp": "2025-06-11T18:39:26.317Z",
      "control_graph_id": "ctg:5388c592189444ad9e84df071c8f3954:8592364792",
      "crawled_timestamp": "2025-06-11T19:39:28.826476115Z",
      "created_timestamp": "2025-06-11T18:40:29.064309181Z",
      "data_domains": ["Endpoint"],
      "description": "A high level detection was triggered on this process for testing purposes.",
      "device": {
        "agent_load_flags": "1",
        "agent_local_time": "2025-06-10T14:00:03.712Z",
        "agent_version": "7.24.19607.0",
        "bios_manufacturer": "Microsoft Corporation",
        "bios_version": "Hyper-V UEFI Release v4.1",
        "cid": "8693deb4bf134cfb8855ee118d9a0243",
        "config_id_base": "65994767",
        "config_id_build": "19607",
        "config_id_platform": "3",
        "device_id": "5388c592189444ad9e84df071c8f3954",
        "external_ip": "20.83.184.209",
        "first_seen": "2025-06-10T01:49:09Z",
        "groups": ["2a8b900d486e4e9eaa024723d6f3742a"],
        "hostinfo": {
          "active_directory_dn_display": ["Domain Controllers"],
          "domain": "normcorp.ai"
        },
        "hostname": "SGNL-CRWD-Proto",
        "instance_id": "3bfaa67d-0dbd-4d49-8a0a-1cb2b0d2e1af",
        "last_seen": "2025-06-11T19:36:15Z",
        "local_ip": "10.3.0.4",
        "mac_address": "00-0d-3a-56-11-c1",
        "machine_domain": "normcorp.ai",
        "major_version": "10",
        "minor_version": "0",
        "modified_timestamp": "2025-06-11T19:38:14Z",
        "os_version": "Windows Server 2022",
        "ou": ["Domain Controllers"],
        "platform_id": "0",
        "platform_name": "Windows",
        "product_type": "2",
        "product_type_desc": "Domain Controller",
        "service_provider": "AZURE",
        "service_provider_account_id": "9405b466-dc55-4b34-a424-2059ff303a68",
        "site_name": "Default-First-Site-Name",
        "status": "normal",
        "system_manufacturer": "Microsoft Corporation",
        "system_product_name": "Virtual Machine"
      },
      "display_name": "TestTriggerHigh",
      "email_sent": true,
      "falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10653769300-10304-81908752?_cid=g04000s6h3lrs7encahh67joerwbbsje",
      "filename": "cmd.exe",
      "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
      "global_prevalence": "common",
      "grandparent_details": {
        "cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
        "filename": "cmd.exe",
        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
        "local_process_id": "7744",
        "md5": "22cdd5b627bed769783a7928efed69ae",
        "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10465963482",
        "process_id": "10465963482",
        "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
        "timestamp": "2025-06-11T17:15:02Z",
        "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_name": "user6"
      },
      "id": "ind:5388c592189444ad9e84df071c8f3954:10653769300-10304-81908752",
      "indicator_id": "ind:5388c592189444ad9e84df071c8f3954:10653769300-10304-81908752",
      "ioc_context": [],
      "local_prevalence": "low",
      "local_process_id": "8000",
      "logon_domain": "normcorp",
      "md5": "22cdd5b627bed769783a7928efed69ae",
      "mitre_attack": [
        {
          "pattern_id": 10304,
          "tactic_id": "CSTA0006",
          "technique_id": "CST0002",
          "tactic": "Falcon Overwatch",
          "technique": "Malicious Activity"
        }
      ],
      "name": "DemoHighPattern",
      "objective": "Falcon Detection Method",
      "parent_details": {
        "cmdline": "cmd  crowdstrike_test_high",
        "filename": "cmd.exe",
        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
        "local_process_id": "8868",
        "md5": "22cdd5b627bed769783a7928efed69ae",
        "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10557972040",
        "process_id": "10557972040",
        "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
        "timestamp": "2025-06-11T17:13:16Z",
        "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_name": "user6"
      },
      "parent_process_id": "10557972040",
      "pattern_disposition": 0,
      "pattern_disposition_description": "Detection, standard detection.",
      "pattern_disposition_details": {
        "blocking_unsupported_or_disabled": false,
        "bootup_safeguard_enabled": false,
        "containment_file_system": false,
        "critical_process_disabled": false,
        "detect": false,
        "fs_operation_blocked": false,
        "handle_operation_downgraded": false,
        "inddet_mask": false,
        "indicator": false,
        "kill_action_failed": false,
        "kill_parent": false,
        "kill_process": false,
        "kill_subprocess": false,
        "mfa_required": false,
        "operation_blocked": false,
        "policy_disabled": false,
        "prevention_provisioning_enabled": false,
        "process_blocked": false,
        "quarantine_file": false,
        "quarantine_machine": false,
        "registry_operation_blocked": false,
        "response_action_already_applied": false,
        "response_action_failed": false,
        "response_action_triggered": false,
        "rooting": false,
        "sensor_only": false,
        "suspend_parent": false,
        "suspend_process": false
      },
      "pattern_id": 10304,
      "platform": "Windows",
      "poly_id": "AACGk960vxNM-4hV7hGNmgJDHT6rPr9ZY7tVWH9KAC5SPAAATiGvJPYbMEpy79ExehQ0j-pi2dGo2QkQs0W4JkPzDurcDg==",
      "priority_explanation": [
        "[MOD] The detection is based on Pattern 10304: A high level detection was triggered on this process for testing purposes.",
        "[MOD] The disposition for the detection: no preventative action was taken",
        "[MOD] The parent process was identified as: cmd.exe"
      ],
      "priority_value": 10,
      "process_id": "10653769300",
      "process_start_time": "1749667166",
      "product": "epp",
      "resolution": "ignored",
      "scenario": "suspicious_activity",
      "seconds_to_resolved": 0,
      "seconds_to_triaged": 0,
      "severity": 70,
      "severity_name": "High",
      "sha1": "0000000000000000000000000000000000000000",
      "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
      "show_in_ui": true,
      "source_products": ["Falcon Insight"],
      "source_vendors": ["CrowdStrike"],
      "status": "closed",
      "tactic": "Falcon Overwatch",
      "tactic_id": "CSTA0006",
      "tags": ["", "ignored"],
      "technique": "Malicious Activity",
      "technique_id": "CST0002",
      "template_instance_id": "1342",
      "timestamp": "2025-06-11T18:39:26.881Z",
      "tree_id": "8592364792",
      "tree_root": "10557972040",
      "triggering_process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10653769300",
      "type": "ldt",
      "updated_timestamp": "2025-06-18T20:23:35.340932619Z",
      "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
      "user_name": "user6",
      "user_principal": "user6@normcorp.ai",
      "incident": {
        "created": "2025-06-16T18:46:54Z",
        "end": "2025-06-16T18:46:54Z",
        "id": "inc:5388c592189444ad9e84df071c8f3954:2343cc087cce4683a3c96b49a1e8c865",
        "score": "8.59047945539694",
        "start": "2025-06-16T18:26:55Z"
      }
    },
    {
      "agent_id": "5388c592189444ad9e84df071c8f3954",
      "aggregate_id": "aggind:5388c592189444ad9e84df071c8f3954:8592364792",
      "alleged_filetype": "exe",
      "cid": "8693deb4bf134cfb8855ee118d9a0243",
      "cloud_indicator": "false",
      "cmdline": "cmd  crowdstrike_test_high",
      "composite_id": "8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10557972040-10304-81543952",
      "confidence": 100,
      "context_timestamp": "2025-06-11T17:13:03.154Z",
      "control_graph_id": "ctg:5388c592189444ad9e84df071c8f3954:8592364792",
      "crawled_timestamp": "2025-06-11T18:13:05.294666965Z",
      "created_timestamp": "2025-06-11T17:14:05.265128027Z",
      "data_domains": ["Endpoint"],
      "description": "A high level detection was triggered on this process for testing purposes.",
      "device": {
        "agent_load_flags": "1",
        "agent_local_time": "2025-06-10T14:00:03.712Z",
        "agent_version": "7.24.19607.0",
        "bios_manufacturer": "Microsoft Corporation",
        "bios_version": "Hyper-V UEFI Release v4.1",
        "cid": "8693deb4bf134cfb8855ee118d9a0243",
        "config_id_base": "65994767",
        "config_id_build": "19607",
        "config_id_platform": "3",
        "device_id": "5388c592189444ad9e84df071c8f3954",
        "external_ip": "20.83.184.209",
        "first_seen": "2025-06-10T01:49:09Z",
        "groups": ["2a8b900d486e4e9eaa024723d6f3742a"],
        "hostinfo": {
          "active_directory_dn_display": ["Domain Controllers"],
          "domain": "normcorp.ai"
        },
        "hostname": "SGNL-CRWD-Proto",
        "instance_id": "3bfaa67d-0dbd-4d49-8a0a-1cb2b0d2e1af",
        "last_seen": "2025-06-11T18:00:17Z",
        "local_ip": "10.3.0.4",
        "mac_address": "00-0d-3a-56-11-c1",
        "machine_domain": "normcorp.ai",
        "major_version": "10",
        "minor_version": "0",
        "modified_timestamp": "2025-06-11T18:12:11Z",
        "os_version": "Windows Server 2022",
        "ou": ["Domain Controllers"],
        "platform_id": "0",
        "platform_name": "Windows",
        "product_type": "2",
        "product_type_desc": "Domain Controller",
        "service_provider": "AZURE",
        "service_provider_account_id": "9405b466-dc55-4b34-a424-2059ff303a68",
        "site_name": "Default-First-Site-Name",
        "status": "normal",
        "system_manufacturer": "Microsoft Corporation",
        "system_product_name": "Virtual Machine"
      },
      "display_name": "TestTriggerHigh",
      "email_sent": true,
      "falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10557972040-10304-81543952?_cid=g04000s6h3lrs7encahh67joerwbbsje",
      "filename": "cmd.exe",
      "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
      "global_prevalence": "common",
      "grandparent_details": {
        "cmdline": "C:\\Windows\\Explorer.EXE",
        "filename": "explorer.exe",
        "filepath": "\\Device\\HarddiskVolume4\\Windows\\explorer.exe",
        "local_process_id": "7740",
        "md5": "d5eaf29530d7c7d703fa0e1f45ff47ba",
        "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:8933771842",
        "process_id": "8933771842",
        "sha256": "e939e7a66aa1e5dfdae89dfd7ee314b60e95cba447082e3f17bc936a22bb6bde",
        "timestamp": "2025-06-11T18:09:49Z",
        "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_name": "user6"
      },
      "id": "ind:5388c592189444ad9e84df071c8f3954:10557972040-10304-81543952",
      "indicator_id": "ind:5388c592189444ad9e84df071c8f3954:10557972040-10304-81543952",
      "ioc_context": [],
      "local_prevalence": "low",
      "local_process_id": "8868",
      "logon_domain": "normcorp",
      "md5": "22cdd5b627bed769783a7928efed69ae",
      "mitre_attack": [
        {
          "pattern_id": 10304,
          "tactic_id": "CSTA0006",
          "technique_id": "CST0002",
          "tactic": "Falcon Overwatch",
          "technique": "Malicious Activity"
        }
      ],
      "name": "DemoHighPattern",
      "objective": "Falcon Detection Method",
      "parent_details": {
        "cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
        "filename": "cmd.exe",
        "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
        "local_process_id": "7744",
        "md5": "22cdd5b627bed769783a7928efed69ae",
        "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10465963482",
        "process_id": "10465963482",
        "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
        "timestamp": "2025-06-11T17:15:02Z",
        "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
        "user_name": "user6"
      },
      "parent_process_id": "10465963482",
      "pattern_disposition": 0,
      "pattern_disposition_description": "Detection, standard detection.",
      "pattern_disposition_details": {
        "blocking_unsupported_or_disabled": false,
        "bootup_safeguard_enabled": false,
        "containment_file_system": false,
        "critical_process_disabled": false,
        "detect": false,
        "fs_operation_blocked": false,
        "handle_operation_downgraded": false,
        "inddet_mask": false,
        "indicator": false,
        "kill_action_failed": false,
        "kill_parent": false,
        "kill_process": false,
        "kill_subprocess": false,
        "mfa_required": false,
        "operation_blocked": false,
        "policy_disabled": false,
        "prevention_provisioning_enabled": false,
        "process_blocked": false,
        "quarantine_file": false,
        "quarantine_machine": false,
        "registry_operation_blocked": false,
        "response_action_already_applied": false,
        "response_action_failed": false,
        "response_action_triggered": false,
        "rooting": false,
        "sensor_only": false,
        "suspend_parent": false,
        "suspend_process": false
      },
      "pattern_id": 10304,
      "platform": "Windows",
      "poly_id": "AACGk960vxNM-4hV7hGNmgJD4A99upX2wfL2y8cFb1M_IgAATiFIDOrk-LktPmcljfnMt9o2NxfND4170PJjqSEguTv-tQ==",
      "priority_explanation": [
        "[MOD] The detection is based on Pattern 10304: A high level detection was triggered on this process for testing purposes.",
        "[MOD] The disposition for the detection: no preventative action was taken",
        "[MOD] The parent process was identified as: cmd.exe"
      ],
      "priority_value": 10,
      "process_id": "10557972040",
      "process_start_time": "1749661983",
      "product": "epp",
      "resolution": "ignored",
      "scenario": "suspicious_activity",
      "seconds_to_resolved": 0,
      "seconds_to_triaged": 0,
      "severity": 70,
      "severity_name": "High",
      "sha1": "0000000000000000000000000000000000000000",
      "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
      "show_in_ui": true,
      "source_products": ["Falcon Insight"],
      "source_vendors": ["CrowdStrike"],
      "status": "closed",
      "tactic": "Falcon Overwatch",
      "tactic_id": "CSTA0006",
      "tags": ["", "ignored"],
      "technique": "Malicious Activity",
      "technique_id": "CST0002",
      "template_instance_id": "1342",
      "timestamp": "2025-06-11T17:13:03.734Z",
      "tree_id": "8592364792",
      "tree_root": "10557972040",
      "triggering_process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10557972040",
      "type": "ldt",
      "updated_timestamp": "2025-06-18T20:23:35.340932619Z",
      "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
      "user_name": "user6",
      "user_principal": "user6@normcorp.ai",
      "incident": {
        "created": "2025-06-16T18:46:54Z",
        "end": "2025-06-16T18:46:54Z",
        "id": "inc:5388c592189444ad9e84df071c8f3954:2343cc087cce4683a3c96b49a1e8c866",
        "score": "8.59047945539694",
        "start": "2025-06-16T18:26:55Z"
      }
    }
  ]
}
`

	AlertResponseLastPage = `{
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
                "agent_id": "5388c592189444ad9e84df071c8f3954",
                "aggregate_id": "aggind:5388c592189444ad9e84df071c8f3954:8591071260",
                "alleged_filetype": "exe",
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "cloud_indicator": "false",
                "cmdline": "cmd  crowdstrike_test_critical",
                "composite_id": "8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10230629714-10303-52340240",
                "confidence": 100,
                "context_timestamp": "2025-06-11T12:58:25.721Z",
                "control_graph_id": "ctg:5388c592189444ad9e84df071c8f3954:8591071260",
                "crawled_timestamp": "2025-06-11T13:58:27.941298086Z",
                "created_timestamp": "2025-06-11T12:59:27.950939533Z",
                "data_domains": [
                    "Endpoint"
                ],
                "description": "A critical level detection was triggered on this process for testing purposes.",
                "device": {
                    "agent_load_flags": "1",
                    "agent_local_time": "2025-06-10T14:00:03.712Z",
                    "agent_version": "7.24.19607.0",
                    "bios_manufacturer": "Microsoft Corporation",
                    "bios_version": "Hyper-V UEFI Release v4.1",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "config_id_base": "65994767",
                    "config_id_build": "19607",
                    "config_id_platform": "3",
                    "device_id": "5388c592189444ad9e84df071c8f3954",
                    "external_ip": "20.83.184.209",
                    "first_seen": "2025-06-10T01:49:09Z",
                    "groups": [
                        "2a8b900d486e4e9eaa024723d6f3742a"
                    ],
                    "hostinfo": {
                        "active_directory_dn_display": [
                            "Domain Controllers"
                        ],
                        "domain": "normcorp.ai"
                    },
                    "hostname": "SGNL-CRWD-Proto",
                    "instance_id": "3bfaa67d-0dbd-4d49-8a0a-1cb2b0d2e1af",
                    "last_seen": "2025-06-11T13:00:17Z",
                    "local_ip": "10.3.0.4",
                    "mac_address": "00-0d-3a-56-11-c1",
                    "machine_domain": "normcorp.ai",
                    "major_version": "10",
                    "minor_version": "0",
                    "modified_timestamp": "2025-06-11T13:54:50Z",
                    "os_version": "Windows Server 2022",
                    "ou": [
                        "Domain Controllers"
                    ],
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "2",
                    "product_type_desc": "Domain Controller",
                    "service_provider": "AZURE",
                    "service_provider_account_id": "9405b466-dc55-4b34-a424-2059ff303a68",
                    "site_name": "Default-First-Site-Name",
                    "status": "normal",
                    "system_manufacturer": "Microsoft Corporation",
                    "system_product_name": "Virtual Machine"
                },
                "display_name": "TestTriggerCritical",
                "email_sent": true,
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10230629714-10303-52340240?_cid=g04000s6h3lrs7encahh67joerwbbsje",
                "filename": "cmd.exe",
                "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                "global_prevalence": "common",
                "grandparent_details": {
                    "cmdline": "cmd  crowdstrike_test_critical",
                    "filename": "cmd.exe",
                    "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                    "local_process_id": "4584",
                    "md5": "22cdd5b627bed769783a7928efed69ae",
                    "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:9781782614",
                    "process_id": "9781782614",
                    "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
                    "timestamp": "2025-06-11T03:05:36Z",
                    "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_name": "user6"
                },
                "id": "ind:5388c592189444ad9e84df071c8f3954:10230629714-10303-52340240",
                "indicator_id": "ind:5388c592189444ad9e84df071c8f3954:10230629714-10303-52340240",
                "ioc_context": [],
                "local_prevalence": "low",
                "local_process_id": "6344",
                "logon_domain": "normcorp",
                "md5": "22cdd5b627bed769783a7928efed69ae",
                "mitre_attack": [
                    {
                        "pattern_id": 10303,
                        "tactic_id": "CSTA0006",
                        "technique_id": "CST0002",
                        "tactic": "Falcon Overwatch",
                        "technique": "Malicious Activity"
                    }
                ],
                "name": "DemoCriticalPattern",
                "objective": "Falcon Detection Method",
                "parent_details": {
                    "cmdline": "cmd  crowdstrike_test_high",
                    "filename": "cmd.exe",
                    "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                    "local_process_id": "6880",
                    "md5": "22cdd5b627bed769783a7928efed69ae",
                    "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10208107226",
                    "process_id": "10208107226",
                    "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
                    "timestamp": "2025-06-11T12:39:56Z",
                    "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_name": "user6"
                },
                "parent_process_id": "10208107226",
                "pattern_disposition": 0,
                "pattern_disposition_description": "Detection, standard detection.",
                "pattern_disposition_details": {
                    "blocking_unsupported_or_disabled": false,
                    "bootup_safeguard_enabled": false,
                    "containment_file_system": false,
                    "critical_process_disabled": false,
                    "detect": false,
                    "fs_operation_blocked": false,
                    "handle_operation_downgraded": false,
                    "inddet_mask": false,
                    "indicator": false,
                    "kill_action_failed": false,
                    "kill_parent": false,
                    "kill_process": false,
                    "kill_subprocess": false,
                    "mfa_required": false,
                    "operation_blocked": false,
                    "policy_disabled": false,
                    "prevention_provisioning_enabled": false,
                    "process_blocked": false,
                    "quarantine_file": false,
                    "quarantine_machine": false,
                    "registry_operation_blocked": false,
                    "response_action_already_applied": false,
                    "response_action_failed": false,
                    "response_action_triggered": false,
                    "rooting": false,
                    "sensor_only": false,
                    "suspend_parent": false,
                    "suspend_process": false
                },
                "pattern_id": 10303,
                "platform": "Windows",
                "poly_id": "AACGk960vxNM-4hV7hGNmgJDUtB2ek1C9RLLaLDRC_hZRQAATiF9rOpMFBt5RWVI-qoMwwpDwkpragqZs7tzcxm_4PlnDQ==",
                "priority_explanation": [
                    "[MOD] The detection is based on Pattern 10303: A critical level detection was triggered on this process for testing purposes.",
                    "[MOD] The disposition for the detection: no preventative action was taken",
                    "[MOD] The parent process was identified as: cmd.exe"
                ],
                "priority_value": 10,
                "process_id": "10230629714",
                "process_start_time": "1749646705",
                "product": "epp",
                "resolution": "ignored",
                "scenario": "suspicious_activity",
                "seconds_to_resolved": 7573,
                "seconds_to_triaged": 0,
                "severity": 90,
                "severity_name": "Critical",
                "sha1": "0000000000000000000000000000000000000000",
                "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
                "show_in_ui": true,
                "source_products": [
                    "Falcon Insight"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "status": "closed",
                "tactic": "Falcon Overwatch",
                "tactic_id": "CSTA0006",
                "tags": [
                    "ignored"
                ],
                "technique": "Malicious Activity",
                "technique_id": "CST0002",
                "template_instance_id": "1343",
                "timestamp": "2025-06-11T12:58:26.269Z",
                "tree_id": "8591071260",
                "tree_root": "9781782614",
                "triggering_process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10230629714",
                "type": "ldt",
                "updated_timestamp": "2025-06-11T16:08:05.654536486Z",
                "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
                "user_name": "user6",
                "incident": {
                    "created": "2025-06-16T18:46:54Z",
                    "end": "2025-06-16T18:46:54Z",
                    "id": "inc:5388c592189444ad9e84df071c8f3954:2343cc087cce4683a3c96b49a1e8c895",
                    "score": "8.59047945539694",
                    "start": "2025-06-16T18:26:55Z"
                }
            },
            {
                "agent_id": "5388c592189444ad9e84df071c8f3954",
                "aggregate_id": "aggind:5388c592189444ad9e84df071c8f3954:8591071260",
                "alleged_filetype": "exe",
                "child_process_ids": [
                    "pid:5388c592189444ad9e84df071c8f3954:10230629714"
                ],
                "cid": "8693deb4bf134cfb8855ee118d9a0243",
                "cloud_indicator": "false",
                "cmdline": "cmd  crowdstrike_test_high",
                "composite_id": "8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10208107226-10304-52110864",
                "confidence": 100,
                "context_timestamp": "2025-06-11T12:39:56.399Z",
                "control_graph_id": "ctg:5388c592189444ad9e84df071c8f3954:8591071260",
                "crawled_timestamp": "2025-06-11T13:39:58.354773222Z",
                "created_timestamp": "2025-06-11T12:40:58.360279214Z",
                "data_domains": [
                    "Endpoint"
                ],
                "description": "A high level detection was triggered on this process for testing purposes.",
                "device": {
                    "agent_load_flags": "1",
                    "agent_local_time": "2025-06-10T14:00:03.712Z",
                    "agent_version": "7.24.19607.0",
                    "bios_manufacturer": "Microsoft Corporation",
                    "bios_version": "Hyper-V UEFI Release v4.1",
                    "cid": "8693deb4bf134cfb8855ee118d9a0243",
                    "config_id_base": "65994767",
                    "config_id_build": "19607",
                    "config_id_platform": "3",
                    "device_id": "5388c592189444ad9e84df071c8f3954",
                    "external_ip": "20.83.184.209",
                    "first_seen": "2025-06-10T01:49:09Z",
                    "groups": [
                        "2a8b900d486e4e9eaa024723d6f3742a"
                    ],
                    "hostinfo": {
                        "active_directory_dn_display": [
                            "Domain Controllers"
                        ],
                        "domain": "normcorp.ai"
                    },
                    "hostname": "SGNL-CRWD-Proto",
                    "instance_id": "3bfaa67d-0dbd-4d49-8a0a-1cb2b0d2e1af",
                    "last_seen": "2025-06-11T13:00:17Z",
                    "local_ip": "10.3.0.4",
                    "mac_address": "00-0d-3a-56-11-c1",
                    "machine_domain": "normcorp.ai",
                    "major_version": "10",
                    "minor_version": "0",
                    "modified_timestamp": "2025-06-11T13:36:06Z",
                    "os_version": "Windows Server 2022",
                    "ou": [
                        "Domain Controllers"
                    ],
                    "platform_id": "0",
                    "platform_name": "Windows",
                    "product_type": "2",
                    "product_type_desc": "Domain Controller",
                    "service_provider": "AZURE",
                    "service_provider_account_id": "9405b466-dc55-4b34-a424-2059ff303a68",
                    "site_name": "Default-First-Site-Name",
                    "status": "normal",
                    "system_manufacturer": "Microsoft Corporation",
                    "system_product_name": "Virtual Machine"
                },
                "display_name": "TestTriggerHigh",
                "email_sent": true,
                "falcon_host_link": "https://falcon.us-2.crowdstrike.com/activity-v2/detections/8693deb4bf134cfb8855ee118d9a0243:ind:5388c592189444ad9e84df071c8f3954:10208107226-10304-52110864?_cid=g04000s6h3lrs7encahh67joerwbbsje",
                "filename": "cmd.exe",
                "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                "global_prevalence": "common",
                "grandparent_details": {
                    "cmdline": "\"C:\\Windows\\system32\\cmd.exe\" ",
                    "filename": "cmd.exe",
                    "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                    "local_process_id": "4468",
                    "md5": "22cdd5b627bed769783a7928efed69ae",
                    "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:9770699358",
                    "process_id": "9770699358",
                    "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
                    "timestamp": "2025-06-11T03:06:36Z",
                    "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_name": "user6"
                },
                "id": "ind:5388c592189444ad9e84df071c8f3954:10208107226-10304-52110864",
                "indicator_id": "ind:5388c592189444ad9e84df071c8f3954:10208107226-10304-52110864",
                "ioc_context": [],
                "local_prevalence": "low",
                "local_process_id": "6880",
                "logon_domain": "normcorp",
                "md5": "22cdd5b627bed769783a7928efed69ae",
                "mitre_attack": [
                    {
                        "pattern_id": 10304,
                        "tactic_id": "CSTA0006",
                        "technique_id": "CST0002",
                        "tactic": "Falcon Overwatch",
                        "technique": "Malicious Activity"
                    }
                ],
                "name": "DemoHighPattern",
                "objective": "Falcon Detection Method",
                "parent_details": {
                    "cmdline": "cmd  crowdstrike_test_critical",
                    "filename": "cmd.exe",
                    "filepath": "\\Device\\HarddiskVolume4\\Windows\\System32\\cmd.exe",
                    "local_process_id": "4584",
                    "md5": "22cdd5b627bed769783a7928efed69ae",
                    "process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:9781782614",
                    "process_id": "9781782614",
                    "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
                    "timestamp": "2025-06-11T03:05:36Z",
                    "user_graph_id": "uid:5388c592189444ad9e84df071c8f3954:S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
                    "user_name": "user6"
                },
                "parent_process_id": "9781782614",
                "pattern_disposition": 0,
                "pattern_disposition_description": "Detection, standard detection.",
                "pattern_disposition_details": {
                    "blocking_unsupported_or_disabled": false,
                    "bootup_safeguard_enabled": false,
                    "containment_file_system": false,
                    "critical_process_disabled": false,
                    "detect": false,
                    "fs_operation_blocked": false,
                    "handle_operation_downgraded": false,
                    "inddet_mask": false,
                    "indicator": false,
                    "kill_action_failed": false,
                    "kill_parent": false,
                    "kill_process": false,
                    "kill_subprocess": false,
                    "mfa_required": false,
                    "operation_blocked": false,
                    "policy_disabled": false,
                    "prevention_provisioning_enabled": false,
                    "process_blocked": false,
                    "quarantine_file": false,
                    "quarantine_machine": false,
                    "registry_operation_blocked": false,
                    "response_action_already_applied": false,
                    "response_action_failed": false,
                    "response_action_triggered": false,
                    "rooting": false,
                    "sensor_only": false,
                    "suspend_parent": false,
                    "suspend_process": false
                },
                "pattern_id": 10304,
                "platform": "Windows",
                "poly_id": "AACGk960vxNM-4hV7hGNmgJDdUpaVtDy_D2QDjhaJeLQWwAATiE6fBKVPmYoRXI_ArwBiCnbF2zzntSdA9OHNKWVWTkQ-w==",
                "priority_explanation": [
                    "[MOD] The detection is based on Pattern 10304: A high level detection was triggered on this process for testing purposes.",
                    "[MOD] The disposition for the detection: no preventative action was taken",
                    "[MOD] The parent process was identified as: cmd.exe"
                ],
                "priority_value": 10,
                "process_id": "10208107226",
                "process_start_time": "1749645596",
                "product": "epp",
                "resolution": "ignored",
                "scenario": "suspicious_activity",
                "seconds_to_resolved": 0,
                "seconds_to_triaged": 0,
                "severity": 70,
                "severity_name": "High",
                "sha1": "0000000000000000000000000000000000000000",
                "sha256": "fc45691495eec0f9de50ccd574ead3a2d83e3088f1a3ad514f46ed84a70efa21",
                "show_in_ui": true,
                "source_products": [
                    "Falcon Insight"
                ],
                "source_vendors": [
                    "CrowdStrike"
                ],
                "status": "closed",
                "tactic": "Falcon Overwatch",
                "tactic_id": "CSTA0006",
                "tags": [
                    "ignored"
                ],
                "technique": "Malicious Activity",
                "technique_id": "CST0002",
                "template_instance_id": "1342",
                "timestamp": "2025-06-11T12:39:56.948Z",
                "tree_id": "8591071260",
                "tree_root": "9781782614",
                "triggering_process_graph_id": "pid:5388c592189444ad9e84df071c8f3954:10208107226",
                "type": "ldt",
                "updated_timestamp": "2025-06-11T16:08:05.654536486Z",
                "user_id": "S-1-5-21-2931850618-1476300705-2742956860-1104",
                "user_name": "user6",
                "incident": {
                    "created": "2025-06-16T18:46:54Z",
                    "end": "2025-06-16T18:46:54Z",
                    "id": "inc:5388c592189444ad9e84df071c8f3954:2343cc087cce4683a3c96b49a1e8c867",
                    "score": "8.59047945539694",
                    "start": "2025-06-16T18:26:55Z"
                }
            }
		]
	}
`

	AlertResponseSpecializedErr = `{
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

	// ************************ Endpoint Incident Responses ************************.
	EndpointIncidentEmptyListResponse = `{
		"meta": {
			"query_time": 0.003925287,
			"pagination": {
				"offset": 0,
				"limit": 100,
				"total": 0
			},
			"powered_by": "incident-api",
			"trace_id": "15c77695-57f4-4fb9-bce2-93a3bfa4dc49"
		},
		"resources": [],
		"errors": []
	}`

	EndpointIncidentEmptyIDsErrorResponse = `{
		"meta": {
			"query_time": 0.000824576,
			"powered_by": "incident-api",
			"trace_id": "3627e1c1-5aa5-433e-bb33-e322d95c1a38"
		},
		"resources": [],
		"errors": [
			{
				"code": 400,
				"message": "The 'ids' parameter must be present at least once and can be present up to 500 times."
			}
		]
	}`

	EndpointIncidentValidResponse = `{
		"meta": {
			"query_time": 0.002345678,
			"powered_by": "incident-api",
			"trace_id": "abc12345-def6-7890-gh12-ijklmn345678"
		},
		"resources": [
			{
				"incident_id": "inc:test123:abc456",
				"state": "closed",
				"status": 20,
				"incident_type": 1,
				"cid": "test-customer-id",
				"host_ids": ["host123", "host456"],
				"created": "2025-01-01T10:00:00Z",
				"start": "2025-01-01T09:30:00Z",
				"end": "2025-01-01T11:00:00Z",
				"fine_score": 75
			}
		],
		"errors": []
	}`
)
