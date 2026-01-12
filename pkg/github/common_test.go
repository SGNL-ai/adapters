// Copyright 2025 SGNL.ai, Inc.

// nolint: lll, goconst
package github_test

import (
	"io"
	"net/http"

	"github.com/sgnl-ai/adapters/pkg/testutil"
)

// Define the endpoints and responses for the mock GitHub server.
// This handler is intended to be re-used throughout the test package.
var TestServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer Testtoken" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{
			"errorCode": "E0000011",
			"errorSummary": "Invalid token provided",
			"errorCursor": "E0000011",
			"errorId": "oaefW5oDjyLRLKVkrmTlp0Thg",
			"errorCauses": []
		}}`))
	}

	query, _ := io.ReadAll(r.Body)
	queryStr := string(query)

	switch r.URL.RequestURI() {

	case "/graphql", "/api/graphql":
		// GraphQL Endpoints
		switch queryStr {
		// Organizations Page 1
		case ValidationQueryBuilder("Organization", "SGNL", 1, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"databaseId": 9,
									"email": null,
									"login": "ArvindOrg1",
									"viewerIsAMember": true,
									"viewerCanCreateTeams": true,
									"updatedAt": "2024-02-02T23:20:22Z",
									"createdAt": "2024-02-02T23:20:22Z"
								}
							]
						}
					}
				}
			}`))
		// Organizations Page 2
		case ValidationQueryBuilder("Organization", "SGNL", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEw",
									"databaseId": 10,
									"email": null,
									"login": "ArvindOrg2",
									"viewerIsAMember": true,
									"viewerCanCreateTeams": true,
									"updatedAt": "2024-02-15T17:00:12Z",
									"createdAt": "2024-02-15T17:00:12Z"
								}
							]
						}
					}
				}
			}`))
		// Organizations Page 3
		case ValidationQueryBuilder("Organization", "SGNL", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": false,
								"endCursor": "Y3Vyc29yOnYyOpKzRW50ZXJwcmlzZVNlcnZlck9yZwU="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"databaseId": 5,
									"email": null,
									"login": "EnterpriseServerOrg",
									"viewerIsAMember": true,
									"viewerCanCreateTeams": true,
									"updatedAt": "2024-01-28T23:00:00Z",
									"createdAt": "2024-01-28T22:59:59Z"
								}
							]
						}
					}
				}
			}`))
		// Teams Page 1 (Includes TeamMembers and TeamRepositories)
		case ValidationQueryBuilder("Team", "SGNL", 2, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"teams": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpMCpXRlYW0xAQ==",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDQ6VGVhbTI=",
												"databaseId": 2,
												"slug": "secret-team-1",
												"viewerCanAdminister": true,
												"updatedAt": "2024-02-02T23:21:54Z",
												"createdAt": "2024-02-02T23:21:54Z",
												"members": {
													"edges": [
														{
															"role": "MAINTAINER",
															"node": {
																"id": "MDQ6VXNlcjY=",
																"databaseId": 6,
																"email": "",
																"login": "arvind",
																"isViewer": true,
																"updatedAt": "2024-01-31T05:09:26Z",
																"createdAt": "2024-01-28T23:28:03Z"
															}
														}
													]
												},
												"repositories": {
													"edges": [
														{
															"permission": "ADMIN",
															"node": {
																"id": "MDEwOlJlcG9zaXRvcnk2",
																"name": "arvindrepo2",
																"databaseId": 6,
																"url": "https://ghe-test-server/ArvindOrg1/arvindrepo2",
																"allowUpdateBranch": false,
																"pushedAt": "2024-02-02T23:22:33Z",
																"createdAt": "2024-02-02T23:22:32Z"
															}
														}
													]
												}
											},
											{
												"id": "MDQ6VGVhbTE=",
												"databaseId": 1,
												"slug": "team1",
												"viewerCanAdminister": true,
												"updatedAt": "2024-02-02T23:21:02Z",
												"createdAt": "2024-02-02T23:21:02Z",
												"members": {
													"edges": [
														{
															"role": "MEMBER",
															"node": {
																"id": "MDQ6VXNlcjQ=",
																"databaseId": 4,
																"email": "",
																"login": "isabella",
																"isViewer": false,
																"updatedAt": "2024-02-22T18:43:44Z",
																"createdAt": "2024-01-28T22:02:26Z"
															}
														},
														{
															"role": "MAINTAINER",
															"node": {
																"id": "MDQ6VXNlcjY=",
																"databaseId": 6,
																"email": "",
																"login": "arvind",
																"isViewer": true,
																"updatedAt": "2024-01-31T05:09:26Z",
																"createdAt": "2024-01-28T23:28:03Z"
															}
														}
													]
												},
												"repositories": {
													"edges": [
														{
															"permission": "MAINTAIN",
															"node": {
																"id": "MDEwOlJlcG9zaXRvcnk1",
																"name": "arvindrepo1",
																"databaseId": 5,
																"url": "https://ghe-test-server/ArvindOrg1/arvindrepo1",
																"allowUpdateBranch": false,
																"pushedAt": "2024-02-02T23:22:20Z",
																"createdAt": "2024-02-02T23:22:20Z"
															}
														},
														{
															"permission": "WRITE",
															"node": {
																"id": "MDEwOlJlcG9zaXRvcnk2",
																"name": "arvindrepo2",
																"databaseId": 6,
																"url": "https://ghe-test-server/ArvindOrg1/arvindrepo2",
																"allowUpdateBranch": false,
																"pushedAt": "2024-02-02T23:22:33Z",
																"createdAt": "2024-02-02T23:22:32Z"
															}
														},
														{
															"permission": "READ",
															"node": {
																"id": "MDEwOlJlcG9zaXRvcnk3",
																"name": "arvindrepo3",
																"databaseId": 7,
																"url": "https://ghe-test-server/ArvindOrg1/arvindrepo3",
																"allowUpdateBranch": false,
																"pushedAt": "2024-02-02T23:22:45Z",
																"createdAt": "2024-02-02T23:22:45Z"
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Teams Page 2 (Includes TeamMembers and TeamRepositories)
		case ValidationQueryBuilder("Team", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEw",
									"teams": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// Teams Page 3 (Includes TeamMembers and TeamRepositories)
		case ValidationQueryBuilder("Team", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKzRW50ZXJwcmlzZVNlcnZlck9yZwU=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"teams": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpMCrXJhbmRvbS10ZWFtLTED",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDQ6VGVhbTM=",
												"databaseId": 3,
												"slug": "random-team-1",
												"viewerCanAdminister": true,
												"updatedAt": "2024-02-16T04:26:06Z",
												"createdAt": "2024-02-16T04:26:06Z",
												"members": {
													"edges": [
														{
															"role": "MAINTAINER",
															"node": {
																"id": "MDQ6VXNlcjY=",
																"databaseId": 6,
																"email": "",
																"login": "arvind",
																"isViewer": true,
																"updatedAt": "2024-01-31T05:09:26Z",
																"createdAt": "2024-01-28T23:28:03Z"
															}
														}
													]
												},
												"repositories": {
													"edges": [
														{
															"permission": "MAINTAIN",
															"node": {
																"id": "MDEwOlJlcG9zaXRvcnkx",
																"name": "enterprise_repo1",
																"databaseId": 1,
																"url": "https://ghe-test-server/EnterpriseServerOrg/enterprise_repo1",
																"allowUpdateBranch": false,
																"pushedAt": "2024-02-02T23:17:27Z",
																"createdAt": "2024-02-02T23:17:26Z"
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Repositories Page 1
		case ValidationQueryBuilder("Repository", "SGNL", 2, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEG",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnk1",
												"name": "arvindrepo1",
												"databaseId": 5,
												"allowUpdateBranch": false,
												"pushedAt": "2024-02-02T23:22:20Z",
												"createdAt": "2024-02-02T23:22:20Z",
												"collaborators": {
													"edges": [
														{
															"permission": "ADMIN",
															"node": {
																"id": "MDQ6VXNlcjQ="
															}
														},
														{
															"permission": "MAINTAIN",
															"node": {
																"id": "MDQ6VXNlcjY="
															}
														}
													]
												}
											},
											{
												"id": "MDEwOlJlcG9zaXRvcnk2",
												"name": "arvindrepo2",
												"databaseId": 6,
												"allowUpdateBranch": false,
												"pushedAt": "2024-02-02T23:22:33Z",
												"createdAt": "2024-02-02T23:22:32Z",
												"collaborators": {
													"edges": [
														{
															"permission": "ADMIN",
															"node": {
																"id": "MDQ6VXNlcjQ="
															}
														}
													]
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=",
								"hasNextPage": true
							}
						}
					}
				}
			}`))
		// Repositories Page 2
		case ValidationQueryBuilder("Repository", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEG")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjk=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEH",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnk3",
												"name": "arvindrepo3",
												"databaseId": 7,
												"allowUpdateBranch": false,
												"pushedAt": "2024-02-02T23:22:45Z",
												"createdAt": "2024-02-02T23:22:45Z",
												"collaborators": {
													"edges": []
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=",
								"hasNextPage": true
							}
						}
					}
				}
			}`))
		// Repositories Page 3
		case ValidationQueryBuilder("Repository", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQk=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEw",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=",
								"hasNextPage": true
							}
						}
					}
				}
			}`))
		// Repositories Page 4
		case ValidationQueryBuilder("Repository", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"name": "enterprise_repo1",
												"databaseId": 1,
												"allowUpdateBranch": false,
												"pushedAt": "2024-02-02T23:17:27Z",
												"createdAt": "2024-02-02T23:17:26Z",
												"collaborators": {
													"edges": []
												}
											},
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"name": "enterprise_repo2",
												"databaseId": 2,
												"allowUpdateBranch": false,
												"pushedAt": "2024-02-02T23:17:42Z",
												"createdAt": "2024-02-02T23:17:41Z",
												"collaborators": {
													"edges": [
														{
															"permission": "MAINTAIN",
															"node": {
																"id": "MDQ6VXNlcjY="
															}
														}
													]
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKzRW50ZXJwcmlzZVNlcnZlck9yZwU=",
								"hasNextPage": false
							}
						}
					}
				}
			}`))
		// Repositories Page 5
		case ValidationQueryBuilder("Repository", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgo="), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpED",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkz",
												"name": "enterprise_repo3",
												"databaseId": 3,
												"allowUpdateBranch": false,
												"pushedAt": "2024-02-02T23:18:01Z",
												"createdAt": "2024-02-02T23:18:01Z",
												"collaborators": {
													"edges": []
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKzRW50ZXJwcmlzZVNlcnZlck9yZwU=",
								"hasNextPage": false
							}
						}
					}
				}
			}`))
		// Users Page 1
		case ValidationQueryBuilder("User", "SGNL", 1, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"membersWithRole": {
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "Y3Vyc29yOnYyOpEE"
										},
										"nodes": [
											{
												"id": "MDQ6VXNlcjQ=",
												"databaseId": 4,
												"email": "",
												"login": "arooxa",
												"isViewer": true,
												"updatedAt": "2024-03-08T04:18:47Z",
												"createdAt": "2024-03-08T04:18:47Z"
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Users Page 2
		case ValidationQueryBuilder("User", "SGNL", 1, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEE")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": true,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"membersWithRole": {
										"pageInfo": {
											"hasNextPage": false,
											"endCursor": "Y3Vyc29yOnYyOpEJ"
										},
										"nodes": [
											{
												"id": "MDQ6VXNlcjk=",
												"databaseId": 9,
												"email": "",
												"login": "isabella-sgnl",
												"isViewer": false,
												"updatedAt": "2024-03-08T19:28:13Z",
												"createdAt": "2024-03-08T17:52:21Z"
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Users Page 3
		case ValidationQueryBuilder("User", "SGNL", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": false,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"membersWithRole": {
										"pageInfo": {
											"hasNextPage": true,
											"endCursor": "Y3Vyc29yOnYyOpEE"
										},
										"nodes": [
											{
												"id": "MDQ6VXNlcjQ=",
												"databaseId": 4,
												"email": "",
												"login": "arooxa",
												"isViewer": true,
												"updatedAt": "2024-03-08T04:18:47Z",
												"createdAt": "2024-03-08T04:18:47Z"
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Users Page 4
		case ValidationQueryBuilder("User", "SGNL", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU="), testutil.GenPtr("Y3Vyc29yOnYyOpEE")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"hasNextPage": false,
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw="
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"membersWithRole": {
										"pageInfo": {
											"hasNextPage": false,
											"endCursor": "Y3Vyc29yOnYyOpEE"
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// Collaborators Page 1
		case ValidationQueryBuilder("Collaborator", "SGNL", 2, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"collaborators": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEJ",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDQ6VXNlcjQ=",
															"databaseId": 4,
															"email": "",
															"login": "arooxa",
															"isViewer": true,
															"updatedAt": "2024-03-08T04:18:47Z",
															"createdAt": "2024-03-08T04:18:47Z"
														},
														{
															"id": "MDQ6VXNlcjk=",
															"databaseId": 9,
															"email": "",
															"login": "isabella-sgnl",
															"isViewer": false,
															"updatedAt": "2024-03-08T19:28:13Z",
															"createdAt": "2024-03-08T17:52:21Z"
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Collaborators Page 2
		case ValidationQueryBuilder("Collaborator", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"collaborators": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEJ",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDQ6VXNlcjQ=",
															"databaseId": 4,
															"email": "",
															"login": "arooxa",
															"isViewer": true,
															"updatedAt": "2024-03-08T04:18:47Z",
															"createdAt": "2024-03-08T04:18:47Z"
														},
														{
															"id": "MDQ6VXNlcjk=",
															"databaseId": 9,
															"email": "",
															"login": "isabella-sgnl",
															"isViewer": false,
															"updatedAt": "2024-03-08T19:28:13Z",
															"createdAt": "2024-03-08T17:52:21Z"
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Collaborators Page 3
		case ValidationQueryBuilder("Collaborator", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEJ")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"collaborators": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEK",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDQ6VXNlcjEw",
															"databaseId": 10,
															"email": "",
															"login": "r-rakshith",
															"isViewer": false,
															"updatedAt": "2024-03-08T17:53:47Z",
															"createdAt": "2024-03-08T17:52:54Z"
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Collaborators Page 4
		case ValidationQueryBuilder("Collaborator", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// OrganizationUsers Page 1
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg1", 1, nil):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
						"membersWithRole": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEE",
								"hasNextPage": true
							},
							"edges": [
								{
									"role": "ADMIN",
									"node": {
										"id": "MDQ6VXNlcjQ=",
										"organizationVerifiedDomainEmails": [
											"arvind@sgnldemos.com"
										]
									}
								}
							]
						}
					}
				}
			}`))
		// OrganizationUsers Page 2
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg1", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")}):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
						"membersWithRole": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEJ",
								"hasNextPage": false
							},
							"edges": [
								{
									"role": "MEMBER",
									"node": {
										"id": "MDQ6VXNlcjk=",
										"organizationVerifiedDomainEmails": [
											"isabella@sgnldemos.com"
										]
									}
								}
							]
						}
					}
				}
			}`))
		// OrganizationUsers Page 3
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg2", 1, nil):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
						"membersWithRole": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEE",
								"hasNextPage": true
							},
							"edges": [
								{
									"role": "ADMIN",
									"node": {
										"id": "MDQ6VXNlcjQ=",
										"organizationVerifiedDomainEmails": []
									}
								}
							]
						}
					}
				}
			}`))
		// OrganizationUsers Page 4
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg2", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")}):
			w.Write([]byte(`{
			"data": {
				"organization": {
					"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
					"membersWithRole": {
						"pageInfo": {
							"endCursor": "Y3Vyc29yOnYyOpEJ",
							"hasNextPage": false
						},
						"edges": []
					}
				}
			}
		}`))
		// Labels Page 1
		case ValidationQueryBuilder("Label", "SGNL", 8, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": "OA",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWwx",
															"name": "bug",
															"color": "d73a4a",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwy",
															"name": "documentation",
															"color": "0075ca",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwz",
															"name": "duplicate",
															"color": "cfd3d7",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWw0",
															"name": "enhancement",
															"color": "a2eeef",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWw1",
															"name": "good first issue",
															"color": "7057ff",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWw2",
															"name": "help wanted",
															"color": "008672",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWw3",
															"name": "invalid",
															"color": "e4e669",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWw4",
															"name": "question",
															"color": "d876e3",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Labels Page 2
		case ValidationQueryBuilder("Label", "SGNL", 8, []*string{nil, nil, testutil.GenPtr("OA")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": "OQ",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWw5",
															"name": "wontfix",
															"color": "ffffff",
															"createdAt": "2024-03-08T18:51:30Z",
															"isDefault": true
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Labels Page 3
		case ValidationQueryBuilder("Label", "SGNL", 8, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"labels": {
													"pageInfo": {
														"endCursor": "OA",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWwxMA==",
															"name": "bug",
															"color": "d73a4a",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxMQ==",
															"name": "documentation",
															"color": "0075ca",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxMg==",
															"name": "duplicate",
															"color": "cfd3d7",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxMw==",
															"name": "enhancement",
															"color": "a2eeef",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxNA==",
															"name": "good first issue",
															"color": "7057ff",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxNQ==",
															"name": "help wanted",
															"color": "008672",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxNg==",
															"name": "invalid",
															"color": "e4e669",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														},
														{
															"id": "MDU6TGFiZWwxNw==",
															"name": "question",
															"color": "d876e3",
															"createdAt": "2024-03-08T18:51:44Z",
															"isDefault": true
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Issues Page 1: Org 1/2, Repo 1/2, Issues [1, 2]
		case ValidationQueryBuilder("Issue", "SGNL", 8, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWUz",
															"title": "issue1",
															"author": {
																"login": "arooxa"
															},
															"createdAt": "2024-03-15T18:40:52Z",
															"isPinned": false
														},
														{
															"id": "MDU6SXNzdWU0",
															"title": "issue2",
															"author": {
																"login": "arooxa"
															},
															"createdAt": "2024-03-15T18:41:04Z",
															"isPinned": false
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Issues Page 2: Org 1/2, Repo 2/2, Issues [1, 2]
		case ValidationQueryBuilder("Issue", "SGNL", 8, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEF",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWUy",
															"title": "issue3",
															"author": {
																"login": "arooxa"
															},
															"createdAt": "2024-03-14T17:43:03Z",
															"isPinned": false
														},
														{
															"id": "MDU6SXNzdWU1",
															"title": "issue4",
															"author": {
																"login": "arooxa"
															},
															"createdAt": "2024-03-15T18:42:01Z",
															"isPinned": false
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// Issues Page 3: Org 2/2 (has no repos)
		case ValidationQueryBuilder("Issue", "SGNL", 8, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// Note: All IssueLabel mock data is from a real instance, however
		// it has been modified because the real data is too sparse.
		// IssueLabels Page 1: Org 1/2, Repo 1/2, Label 1/3, Issue [1, 1]/1
		case ValidationQueryBuilder("IssueLabel", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": "MQ",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWwx",
															"issues": {
																"pageInfo": {
																	"endCursor": "Y3Vyc29yOnYyOpED",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDU6SXNzdWUz",
																		"title": "issue1"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueLabels Page 2: Org 1/2, Repo 1/2, Label 2/3, (has no issues)
		case ValidationQueryBuilder("IssueLabel", "SGNL", 5, []*string{nil, nil, testutil.GenPtr("MQ")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": "Mg",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWwy",
															"issues": {
																"pageInfo": {
																	"endCursor": null,
																	"hasNextPage": false
																},
																"nodes": []
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueLabels Page 3: Org 1/2, Repo 1/2, Label 3/3, Issue [1, 2]/2
		case ValidationQueryBuilder("IssueLabel", "SGNL", 5, []*string{nil, nil, testutil.GenPtr("Mg")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": null,
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWwy",
															"issues": {
																"pageInfo": {
																	"endCursor": null,
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDU6SXNzdWUz",
																		"title": "issue1"
																	},
																	{
																		"id": "MDU6SXNzdWU0",
																		"title": "issue2"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueLabels Page 4: Org 1/2, Repo 2/2, (no labels)
		case ValidationQueryBuilder("IssueLabel", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx2",
												"labels": {
													"pageInfo": {
														"endCursor": null,
														"hasNextPage": false
													},
													"nodes": []
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueLabels Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("IssueLabel", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": null,
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// Note: All PullRequestLabel mock data is from a real instance, however it
		// has been modified because the real data is too sparse.
		// PullRequestLabels Page 1: Org 1/2, Repo 1/2, Label [1]/2, PullRequest [1]/1
		case ValidationQueryBuilder("PullRequestLabel", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": "NQ",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWw0",
															"pullRequests": {
																"pageInfo": {
																	"endCursor": "Y3Vyc29yOnYyOpEB",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
																		"title": "Create README.md"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestLabels Page 2: Org 1/2, Repo 1/2, Label [2]/2, (has no pull requests)
		case ValidationQueryBuilder("PullRequestLabel", "SGNL", 5, []*string{nil, nil, testutil.GenPtr("NQ")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": null,
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWwx",
															"pullRequests": {
																"pageInfo": {
																	"endCursor": null,
																	"hasNextPage": false
																},
																"nodes": []
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestLabels Page 3: Org 1/2, Repo 2/2, Label [1]/1, PullRequest [1, 2]/2
		case ValidationQueryBuilder("PullRequestLabel", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"labels": {
													"pageInfo": {
														"endCursor": null,
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6TGFiZWw1",
															"pullRequests": {
																"pageInfo": {
																	"endCursor": "Y3Vyc29yOnYyOpEB",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDExOlB1bGxSZXF1ZXsdsd$S0=",
																		"title": "BRANCH4PR"
																	},
																	{
																		"id": "MDExOlB1bGxSZXFsssd@@",
																		"title": "BRANCH5PR UPDATE README"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestLabels Page 4: Org 2/2, (has no repos)
		case ValidationQueryBuilder("PullRequestLabel", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// IssueAssignees Page 1: Org 1/2, Repo 1/2, Issue 1/2, Assignees [1, 2]
		case ValidationQueryBuilder("IssueAssignee", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWUz",
															"assignees": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjk=",
																		"login": "isabella-sgnl"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueAssignees Page 2: Org 1/2, Repo 1/2, Issue 2/2, Assignees [1, 1]
		case ValidationQueryBuilder("IssueAssignee", "SGNL", 5, []*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWU0",
															"assignees": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjk=",
																		"login": "isabella-sgnl"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueAssignees Page 3: Org 1/2, Repo 2/2, Issue 1/2, Assignees [1, 2]
		case ValidationQueryBuilder("IssueAssignee", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWUy",
															"assignees": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjk=",
																		"login": "isabella-sgnl"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueAssignees Page 4: Org 1/2, Repo 2/2, Issue 2/2, Assignees [1, 2]
		case ValidationQueryBuilder("IssueAssignee", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEF",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWU1",
															"assignees": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjEw",
																		"login": "r-rakshith"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueAssignees Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("IssueAssignee", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// IssueParticipants Page 1: Org 1/2, Repo 1/2, Issue 1/2, Participants [1, 1]
		case ValidationQueryBuilder("IssueParticipant", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"name": "repo1",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWUz",
															"title": "issue1",
															"participants": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueParticipants Page 2: Org 1/2, Repo 1/2, Issue 2/2, Participants [1, 1]
		case ValidationQueryBuilder("IssueParticipant", "SGNL", 5, []*string{nil, nil, testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"name": "repo1",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWU0",
															"title": "issue2",
															"participants": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueParticipants Page 3: Org 1/2, Repo 2/2, Issue 1/2, Participants [1, 2]
		case ValidationQueryBuilder("IssueParticipant", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"name": "repo2",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWUy",
															"title": "issue3",
															"participants": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjEw",
																		"login": "r-rakshith"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueParticipants Page 4: Org 1/2, Repo 2/2, Issue 2/2, Participants [1, 2]
		case ValidationQueryBuilder("IssueParticipant", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"name": "repo2",
												"issues": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEF",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDU6SXNzdWU1",
															"title": "issue4",
															"participants": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjEw",
																		"login": "r-rakshith"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// IssueParticipants Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("IssueParticipant", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequests Page 1: Org 1/2, Repo 1/2, PR 1/1
		case ValidationQueryBuilder("PullRequest", "SGNL", 2, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"title": "Create README.md",
															"closed": false,
															"createdAt": "2024-03-13T23:07:49Z",
															"author": {
																"login": "arooxa"
															},
															"baseRepository": {
																"id": "MDEwOlJlcG9zaXRvcnkx"
															},
															"headRepository": {
																"id": "MDEwOlJlcG9zaXRvcnkx"
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequests Page 2: Org 1/2, Repo 2/2, PR [1, 2]/3
		case ValidationQueryBuilder("PullRequest", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mg==",
															"title": "[branch4PR] README",
															"closed": false,
															"createdAt": "2024-03-15T18:43:27Z",
															"author": {
																"login": "arooxa"
															},
															"baseRepository": {
																"id": "MDEwOlJlcG9zaXRvcnky"
															},
															"headRepository": {
																"id": "MDEwOlJlcG9zaXRvcnky"
															}
														},
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mw==",
															"title": "[branch5PR] README.md",
															"closed": false,
															"createdAt": "2024-03-15T18:46:54Z",
															"author": {
																"login": "arooxa"
															},
															"baseRepository": {
																"id": "MDEwOlJlcG9zaXRvcnky"
															},
															"headRepository": {
																"id": "MDEwOlJlcG9zaXRvcnky"
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequests Page 3: Org 1/2, Repo 2/2, PR [3]/3
		case ValidationQueryBuilder("PullRequest", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0NA==",
															"title": "[branch6PR] readMe",
															"closed": false,
															"createdAt": "2024-03-15T22:40:43Z",
															"author": {
																"login": "arooxa"
															},
															"baseRepository": {
																"id": "MDEwOlJlcG9zaXRvcnky"
															},
															"headRepository": {
																"id": "MDEwOlJlcG9zaXRvcnky"
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequests Page 4: Org 2/2 (has no repos)
		case ValidationQueryBuilder("PullRequest", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestChangedFiles Page 1: Org 1/2, Repo 1/2, PR 1/1, Files [1]/1
		case ValidationQueryBuilder("PullRequestChangedFile", "SGNL", 2, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"files": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"path": "README.md",
																		"changeType": "ADDED"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestChangedFiles Page 2: Org 1/2, Repo 2/2, PR 1/3, Files [1]/1
		case ValidationQueryBuilder("PullRequestChangedFile", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mg==",
															"files": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"path": "random/file.txt",
																		"changeType": "DELETED"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestChangedFiles Page 3: Org 1/2, Repo 2/2, PR 2/3, Files [1]/1
		case ValidationQueryBuilder("PullRequestChangedFile", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mw==",
															"files": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"path": "README.md",
																		"changeType": "ADDED"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestChangedFiles Page 4: Org 1/2, Repo 2/2, PR 3/3, Files [1]/1
		case ValidationQueryBuilder("PullRequestChangedFile", "SGNL", 2, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0NA==",
															"files": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"path": "random/file.txt",
																		"changeType": "ADDED"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestChangedFiles Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("PullRequestChangedFile", "SGNL", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestAssignee Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Assignees [1]/1
		case ValidationQueryBuilder("PullRequestAssignee", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"assignees": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestAssignee Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Assignees [1, 2]/2
		case ValidationQueryBuilder("PullRequestAssignee", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mg==",
															"assignees": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjk=",
																		"login": "isabella-sgnl"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestAssignee Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Assignees [1, 1]/1
		case ValidationQueryBuilder("PullRequestAssignee", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mw==",
															"assignees": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestAssignee Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, (has no assignees)
		case ValidationQueryBuilder("PullRequestAssignee", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0NA==",
															"assignees": {
																"pageInfo": {
																	"endCursor": null,
																	"hasNextPage": false
																},
																"nodes": []
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestAssignee Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("PullRequestAssignee", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestParticipant Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Participants [1]/1
		case ValidationQueryBuilder("PullRequestParticipant", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"participants": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestParticipant Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Participants [1, 2]/2
		case ValidationQueryBuilder("PullRequestParticipant", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mg==",
															"participants": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjEw",
																		"login": "r-rakshith"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestParticipant Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Participants [1, 2]/2
		case ValidationQueryBuilder("PullRequestParticipant", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mw==",
															"participants": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	},
																	{
																		"id": "MDQ6VXNlcjEw",
																		"login": "r-rakshith"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestParticipant Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Participants [1]/1
		case ValidationQueryBuilder("PullRequestParticipant", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0NA==",
															"participants": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDQ6VXNlcjQ=",
																		"login": "arooxa"
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestParticipant Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("PullRequestParticipant", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestCommit Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, Commits [1]/1
		case ValidationQueryBuilder("PullRequestCommit", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"commits": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh",
																		"commit": {
																			"id": "MDY6Q29tbWl0MTo0YWNkMDEzNTJkNTZjYTMzMTA1ZmMyMjU4ZDFmMTI4NzZmMzhlZjRh",
																			"committedDate": "2024-03-13T23:07:39Z",
																			"author": {
																				"email": "arvind@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjQ=",
																					"login": "arooxa"
																				}
																			}
																		}
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestCommit Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Commits [1, 3]/3
		case ValidationQueryBuilder("PullRequestCommit", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mg==",
															"commits": {
																"pageInfo": {
																	"endCursor": "Mw",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MjozZjBiMmRiMDM3NmJjYTgwNjM0NDRmNjI4ZWI3ZWI5Y2U4NTk1ZGNj",
																		"commit": {
																			"id": "MDY6Q29tbWl0MjozZjBiMmRiMDM3NmJjYTgwNjM0NDRmNjI4ZWI3ZWI5Y2U4NTk1ZGNj",
																			"committedDate": "2024-03-15T18:43:10Z",
																			"author": {
																				"email": "arvind@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjQ=",
																					"login": "arooxa"
																				}
																			}
																		}
																	},
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0Mjo2MTFlOTU3NGUzODNiNWQ2NmVjNjAwNDMxYTg4ODRkMzc4OGJiMTQx",
																		"commit": {
																			"id": "MDY6Q29tbWl0Mjo2MTFlOTU3NGUzODNiNWQ2NmVjNjAwNDMxYTg4ODRkMzc4OGJiMTQx",
																			"committedDate": "2024-03-16T21:18:12Z",
																			"author": {
																				"email": "arvind@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjQ=",
																					"login": "arooxa"
																				}
																			}
																		}
																	},
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0MjpjMWMzNmQ2ZWQ0M2U4ZmVmMjlhNGExNTc2ZWQxZTYxNGZkMGMzNDFi",
																		"commit": {
																			"id": "MDY6Q29tbWl0MjpjMWMzNmQ2ZWQ0M2U4ZmVmMjlhNGExNTc2ZWQxZTYxNGZkMGMzNDFi",
																			"committedDate": "2024-03-22T21:48:21Z",
																			"author": {
																				"email": "rakshith@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjEw",
																					"login": "r-rakshith"
																				}
																			}
																		}
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestCommit Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Commits [1]/1
		case ValidationQueryBuilder("PullRequestCommit", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mw==",
															"commits": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0Mzo1OTYxZGE3NDk1NmJhNWRiYTQ0YWEyYjQ4Mjc2MzM4MGNkNDhhMWZj",
																		"commit": {
																			"id": "MDY6Q29tbWl0Mjo1OTYxZGE3NDk1NmJhNWRiYTQ0YWEyYjQ4Mjc2MzM4MGNkNDhhMWZj",
																			"committedDate": "2024-03-15T18:45:03Z",
																			"author": {
																				"email": "arvind@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjQ=",
																					"login": "arooxa"
																				}
																			}
																		}
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestCommit Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Commits [1, 2]/2
		case ValidationQueryBuilder("PullRequestCommit", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0NA==",
															"commits": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0NDo1YTVlNzJmNWQwZjk0MjVlZTk3NDc4NzMxZTc2MDczYjBmMTYzY2Fi",
																		"commit": {
																			"id": "MDY6Q29tbWl0Mjo1YTVlNzJmNWQwZjk0MjVlZTk3NDc4NzMxZTc2MDczYjBmMTYzY2Fi",
																			"committedDate": "2024-03-15T22:39:33Z",
																			"author": {
																				"email": "arvind@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjQ=",
																					"login": "arooxa"
																				}
																			}
																		}
																	},
																	{
																		"id": "MDE3OlB1bGxSZXF1ZXN0Q29tbWl0NDpkNjE2NmYwYTlmMmQwMGZlYmFjYzZhYTM3MTAwYWY0YzAxNzBlYzhk",
																		"commit": {
																			"id": "MDY6Q29tbWl0MjpkNjE2NmYwYTlmMmQwMGZlYmFjYzZhYTM3MTAwYWY0YzAxNzBlYzhk",
																			"committedDate": "2024-03-15T22:44:24Z",
																			"author": {
																				"email": "arvind@sgnl.ai",
																				"user": {
																					"id": "MDQ6VXNlcjQ=",
																					"login": "arooxa"
																				}
																			}
																		}
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestCommit Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("PullRequestCommit", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestReview Page 1: Org 1/2, Repo 1/2, PullRequest 1/1, (has no reviews)
		case ValidationQueryBuilder("PullRequestReview", "SGNL", 5, nil):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEB",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnkx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"latestOpinionatedReviews": {
																"pageInfo": {
																	"endCursor": null,
																	"hasNextPage": false
																},
																"nodes": []
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestReview Page 2: Org 1/2, Repo 2/2, PullRequest 1/3, Reviews [1]/1
		case ValidationQueryBuilder("PullRequestReview", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEC",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mg==",
															"latestOpinionatedReviews": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"author": {
																			"login": "r-rakshith"
																		},
																		"state": "APPROVED",
																		"id": "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3NQ==",
																		"createdAt": "2024-03-15T21:05:52Z",
																		"authorCanPushToRepository": true
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestReview Page 3: Org 1/2, Repo 2/2, PullRequest 2/3, Reviews [1, 2]/2
		case ValidationQueryBuilder("PullRequestReview", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpEC")}):
			w.Write([]byte(`{
					"data": {
						"enterprise": {
							"id": "MDEwOkVudGVycHJpc2Ux",
							"organizations": {
								"pageInfo": {
									"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
									"hasNextPage": true
								},
								"nodes": [
									{
										"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpED",
														"hasNextPage": true
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0Mw==",
															"latestOpinionatedReviews": {
																"pageInfo": {
																	"endCursor": "Mg",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"author": {
																			"login": "r-rakshith"
																		},
																		"state": "APPROVED",
																		"id": "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Ng==",
																		"createdAt": "2024-03-15T21:06:25Z",
																		"authorCanPushToRepository": true
																	},
																	{
																		"author": {
																			"login": "isabella-sgnl"
																		},
																		"state": "APPROVED",
																		"id": "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3Mg==",
																		"createdAt": "2024-03-15T19:45:09Z",
																		"authorCanPushToRepository": true
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestReview Page 4: Org 1/2, Repo 2/2, PullRequest 3/3, Reviews [1]/1
		case ValidationQueryBuilder("PullRequestReview", "SGNL", 5, []*string{nil, testutil.GenPtr("Y3Vyc29yOnYyOpEB"), testutil.GenPtr("Y3Vyc29yOnYyOpED")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
									"repositories": {
										"pageInfo": {
											"endCursor": "Y3Vyc29yOnYyOpEC",
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDEwOlJlcG9zaXRvcnky",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEE",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0NA==",
															"latestOpinionatedReviews": {
																"pageInfo": {
																	"endCursor": "MQ",
																	"hasNextPage": false
																},
																"nodes": [
																	{
																		"author": {
																			"login": "isabella-sgnl"
																		},
																		"state": "CHANGES_REQUESTED",
																		"id": "MDE3OlB1bGxSZXF1ZXN0UmV2aWV3OA==",
																		"createdAt": "2024-03-15T22:46:20Z",
																		"authorCanPushToRepository": true
																	}
																]
															}
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestReview Page 5: Org 2/2 (has no repos)
		case ValidationQueryBuilder("PullRequestReview", "SGNL", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMQU=")}):
			w.Write([]byte(`{
				"data": {
					"enterprise": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"organizations": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpKqQXJ2aW5kT3JnMgw=",
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
									"repositories": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": []
									}
								}
							]
						}
					}
				}
			}`))
		// Repository using organization and no eneterprise slug
		case ValidationQueryBuilder("Repository", "arvindorg1", 2, nil, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "O_kgDOCPwuWw",
						"repositories": {
							"nodes": [
								{
									"id": "MDEwOlJlcG9zaXRvcnk1",
									"name": "arvindrepo1",
									"databaseId": 5,
									"allowUpdateBranch": false,
									"pushedAt": "2024-02-02T23:22:20Z",
									"createdAt": "2024-02-02T23:22:20Z",
									"collaborators": {
										"edges": [
											{
												"permission": "ADMIN",
												"node": {
													"id": "MDQ6VXNlcjQ="
												}
											},
											{
												"permission": "MAINTAIN",
												"node": {
													"id": "MDQ6VXNlcjY="
												}
											}
										]
									}
								},
								{
									"id": "MDEwOlJlcG9zaXRvcnk2",
									"name": "arvindrepo2",
									"databaseId": 6,
									"allowUpdateBranch": false,
									"pushedAt": "2024-02-02T23:22:33Z",
									"createdAt": "2024-02-02T23:22:32Z",
									"collaborators": {
										"edges": [
											{
												"permission": "ADMIN",
												"node": {
													"id": "MDQ6VXNlcjQ="
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEG",
								"hasNextPage": true
							}
						}
					}
				}
			}`))
		// Last Page of Repository for organization "arvindorg1" and no eneterprise slug.
		case ValidationQueryBuilder("Repository", "arvindorg1", 2, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpEG")}, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "O_kgDOCPwuWw",
						"repositories": {
							"nodes": [
								{
									"id": "MDEwOlJlcG9zaXRvcnk1",
									"name": "arvindrepo3",
									"databaseId": 7,
									"allowUpdateBranch": false,
									"pushedAt": "2024-02-02T23:22:20Z",
									"createdAt": "2024-02-02T23:22:20Z",
									"collaborators": {
										"edges": [
											{
												"permission": "ADMIN",
												"node": {
													"id": "MDQ6VXNlcjQ="
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOp==",
								"hasNextPage": false
							}
						}
					}
				}
			}`))
		// Last Page of Repository for organization "arvindorg2" and no eneterprise slug.
		case ValidationQueryBuilder("Repository", "arvindorg2", 2, nil, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "O_kgDOCPwuXxyz",
						"repositories": {
							"nodes": [
								{
									"id": "MDEwOlJlcG9zaXRvcnabc",
									"name": "arvindrepo4",
									"databaseId": 9,
									"allowUpdateBranch": false,
									"pushedAt": "2024-02-02T23:22:20Z",
									"createdAt": "2024-02-02T23:22:20Z",
									"collaborators": {
										"edges": [
											{
												"permission": "ADMIN",
												"node": {
													"id": "MDQ6VXNlcjQ="
												}
											}
										]
									}
								}
							],
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpzzz56==",
								"hasNextPage": false
							}
						}
					}
				}
			}`))
		// PullRequestLabels Page 1: Org 1/2, Repo 1/2, Labels 1/2, PullRequest 1/1
		case ValidationQueryBuilder("PullRequestLabel", "arvindorg1", 5, nil, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEwOkVudGVycHJpc2Ux",
						"repositories": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEB",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEwOlJlcG9zaXRvcnkx",
									"labels": {
										"pageInfo": {
											"endCursor": "NQ",
											"hasNextPage": true
										},
										"nodes": [
											{
												"id": "MDU6TGFiZWw0",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXN0MQ==",
															"title": "Create README.md"
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestLabels Page 2: Org 1/2, Repo 1/2, Label 2/2 - has no Pull requests
		case ValidationQueryBuilder("PullRequestLabel", "arvindorg1", 5, []*string{nil, testutil.GenPtr("NQ")}, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEwOkVudGVycHJpc2Uz",
						"repositories": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEB",
								"hasNextPage": true
							},
							"nodes": [
								{
									"id": "MDEwOlJlcG9zaXRvcnkx",
									"labels": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDU6TGFiZWwx",
												"pullRequests": {
													"pageInfo": {
														"endCursor": null,
														"hasNextPage": false
													},
													"nodes": []
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		case ValidationQueryBuilder("PullRequestLabel", "arvindorg1", 5, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpEB")}, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
						"repositories": {
							"pageInfo": {
								"endCursor": null,
								"hasNextPage": false
							},
							"nodes": [
								{
									"id": "MDEwOlJlcG9zaXRvcnkx",
									"labels": {
										"pageInfo": {
											"endCursor": null,
											"hasNextPage": false
										},
										"nodes": [
											{
												"id": "MDU6TGFiZWw1",
												"pullRequests": {
													"pageInfo": {
														"endCursor": "Y3Vyc29yOnYyOpEB",
														"hasNextPage": false
													},
													"nodes": [
														{
															"id": "MDExOlB1bGxSZXF1ZXsdsd$S0=",
															"title": "BRANCH4PR"
														},
														{
															"id": "MDExOlB1bGxSZXFsssd@@",
															"title": "BRANCH5PR UPDATE README"
														}
													]
												}
											}
										]
									}
								}
							]
						}
					}
				}
			}`))
		// PullRequestLabels Page 4: Org 2/2, (has no repos)
		case ValidationQueryBuilder("PullRequestLabel", "arvindorg2", 5, nil, "arvindorg1", "arvindorg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
						"repositories": {
							"pageInfo": {
								"endCursor": null,
								"hasNextPage": false
							},
							"nodes": []
						}
					}
				}
			}`))
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg1", 1, nil, "ArvindOrg1", "ArvindOrg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
						"membersWithRole": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEE",
								"hasNextPage": true
							},
							"edges": [
								{
									"role": "ADMIN",
									"node": {
										"id": "MDQ6VXNlcjQ=",
										"organizationVerifiedDomainEmails": [
											"arvind@sgnldemos.com"
										]
									}
								}
							]
						}
					}
				}
			}`))
		// OrganizationUsers Page 2
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg1", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")}, "ArvindOrg1", "ArvindOrg2"):
			w.Write([]byte(`{
			"data": {
				"organization": {
					"id": "MDEyOk9yZ2FuaXphdGlvbjU=",
					"membersWithRole": {
						"pageInfo": {
							"endCursor": "Y3Vyc29yOnYyOpEJ",
							"hasNextPage": false
						},
						"edges": [
							{
								"role": "MEMBER",
								"node": {
									"id": "MDQ6VXNlcjk=",
									"organizationVerifiedDomainEmails": [
										"isabella@sgnldemos.com"
									]
								}
							}
						]
					}
				}
			}
		}`))
		// OrganizationUsers Page 3 - 2nd Org has only one member.
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg2", 1, nil, "ArvindOrg1", "ArvindOrg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
						"membersWithRole": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEE",
								"hasNextPage": false
							},
							"edges": [
								{
									"role": "ADMIN",
									"node": {
										"id": "MDQ6VXNlcjQ=",
										"organizationVerifiedDomainEmails": []
									}
								}
							]
						}
					}
				}
			}`))
		// OrganizationUsers Page 4 - No more users.
		case ValidationQueryBuilder("OrganizationUser", "ArvindOrg2", 1, []*string{testutil.GenPtr("Y3Vyc29yOnYyOpEE")}, "ArvindOrg1", "ArvindOrg2"):
			w.Write([]byte(`{
				"data": {
					"organization": {
						"id": "MDEyOk9yZ2FuaXphdGlvbjEy",
						"membersWithRole": {
							"pageInfo": {
								"endCursor": "Y3Vyc29yOnYyOpEJ",
								"hasNextPage": false
							},
							"edges": []
						}
					}
				}
			}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Not Found"}`))
		}
	default:
		// REST Endpoints

		if queryStr != "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"message": "Body was not empty for REST request."}`))

			break
		}

		switch r.URL.RequestURI() {
		// Note: All Secret Scanning Alert Responses are Mock Data from the GitHub docs. (not from live instance)
		// https://docs.github.com/en/rest/secret-scanning/secret-scanning?apiVersion=2022-11-28#list-secret-scanning-alerts-for-an-enterprise
		// Secret Scanning Alerts Page 1
		case "/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1":
			w.Header().Add("Link",
				`<https://test-instance.com/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2>; rel="next",
				<https://test-instance.com/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2>; rel="last"`,
			)
			w.Write([]byte(`[
				{
					"number": 2,
					"created_at": "2020-11-06T18:48:51Z",
					"url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2",
					"html_url": "https://github.com/owner/private-repo/security/secret-scanning/2",
					"locations_url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2/locations",
					"state": "resolved",
					"resolution": "false_positive",
					"resolved_at": "2020-11-07T02:47:13Z",
					"resolved_by": {
					  "login": "monalisa",
					  "id": 2,
					  "node_id": "MDQ6VXNlcjI=",
					  "avatar_url": "https://alambic.github.com/avatars/u/2?",
					  "gravatar_id": "",
					  "url": "https://api.github.com/users/monalisa",
					  "html_url": "https://github.com/monalisa",
					  "followers_url": "https://api.github.com/users/monalisa/followers",
					  "following_url": "https://api.github.com/users/monalisa/following{/other_user}",
					  "gists_url": "https://api.github.com/users/monalisa/gists{/gist_id}",
					  "starred_url": "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
					  "subscriptions_url": "https://api.github.com/users/monalisa/subscriptions",
					  "organizations_url": "https://api.github.com/users/monalisa/orgs",
					  "repos_url": "https://api.github.com/users/monalisa/repos",
					  "events_url": "https://api.github.com/users/monalisa/events{/privacy}",
					  "received_events_url": "https://api.github.com/users/monalisa/received_events",
					  "type": "User",
					  "site_admin": true
					},
					"secret_type": "adafruit_io_key",
					"secret_type_display_name": "Adafruit IO Key",
					"secret": "aio_XXXXXXXXXXXXXXXXXXXXXXXXXXXX",
					"repository": {
					  "id": 1296269,
					  "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
					  "name": "Hello-World",
					  "full_name": "octocat/Hello-World",
					  "owner": {
						"login": "octocat",
						"id": 1,
						"node_id": "MDQ6VXNlcjE=",
						"avatar_url": "https://github.com/images/error/octocat_happy.gif",
						"gravatar_id": "",
						"url": "https://api.github.com/users/octocat",
						"html_url": "https://github.com/octocat",
						"followers_url": "https://api.github.com/users/octocat/followers",
						"following_url": "https://api.github.com/users/octocat/following{/other_user}",
						"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
						"organizations_url": "https://api.github.com/users/octocat/orgs",
						"repos_url": "https://api.github.com/users/octocat/repos",
						"events_url": "https://api.github.com/users/octocat/events{/privacy}",
						"received_events_url": "https://api.github.com/users/octocat/received_events",
						"type": "User",
						"site_admin": false
					  },
					  "private": false,
					  "html_url": "https://github.com/octocat/Hello-World",
					  "description": "This your first repo!",
					  "fork": false,
					  "url": "https://api.github.com/repos/octocat/Hello-World",
					  "archive_url": "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
					  "assignees_url": "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
					  "blobs_url": "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
					  "branches_url": "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
					  "collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
					  "comments_url": "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
					  "commits_url": "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
					  "compare_url": "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
					  "contents_url": "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
					  "contributors_url": "https://api.github.com/repos/octocat/Hello-World/contributors",
					  "deployments_url": "https://api.github.com/repos/octocat/Hello-World/deployments",
					  "downloads_url": "https://api.github.com/repos/octocat/Hello-World/downloads",
					  "events_url": "https://api.github.com/repos/octocat/Hello-World/events",
					  "forks_url": "https://api.github.com/repos/octocat/Hello-World/forks",
					  "git_commits_url": "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
					  "git_refs_url": "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
					  "git_tags_url": "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
					  "issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
					  "issue_events_url": "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
					  "issues_url": "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
					  "keys_url": "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
					  "labels_url": "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
					  "languages_url": "https://api.github.com/repos/octocat/Hello-World/languages",
					  "merges_url": "https://api.github.com/repos/octocat/Hello-World/merges",
					  "milestones_url": "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
					  "notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
					  "pulls_url": "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
					  "releases_url": "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
					  "stargazers_url": "https://api.github.com/repos/octocat/Hello-World/stargazers",
					  "statuses_url": "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
					  "subscribers_url": "https://api.github.com/repos/octocat/Hello-World/subscribers",
					  "subscription_url": "https://api.github.com/repos/octocat/Hello-World/subscription",
					  "tags_url": "https://api.github.com/repos/octocat/Hello-World/tags",
					  "teams_url": "https://api.github.com/repos/octocat/Hello-World/teams",
					  "trees_url": "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
					  "hooks_url": "https://api.github.com/repos/octocat/Hello-World/hooks"
					},
					"push_protection_bypassed_by": {
					  "login": "monalisa",
					  "id": 2,
					  "node_id": "MDQ6VXNlcjI=",
					  "avatar_url": "https://alambic.github.com/avatars/u/2?",
					  "gravatar_id": "",
					  "url": "https://api.github.com/users/monalisa",
					  "html_url": "https://github.com/monalisa",
					  "followers_url": "https://api.github.com/users/monalisa/followers",
					  "following_url": "https://api.github.com/users/monalisa/following{/other_user}",
					  "gists_url": "https://api.github.com/users/monalisa/gists{/gist_id}",
					  "starred_url": "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
					  "subscriptions_url": "https://api.github.com/users/monalisa/subscriptions",
					  "organizations_url": "https://api.github.com/users/monalisa/orgs",
					  "repos_url": "https://api.github.com/users/monalisa/repos",
					  "events_url": "https://api.github.com/users/monalisa/events{/privacy}",
					  "received_events_url": "https://api.github.com/users/monalisa/received_events",
					  "type": "User",
					  "site_admin": true
					},
					"push_protection_bypassed": true,
					"push_protection_bypassed_at": "2020-11-06T21:48:51Z",
					"resolution_comment": "Example comment",
					"validity": "active"
				  }
			]`))

		// Secret Scanning Alerts Page 2
		case "/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=2":
			w.Header().Add("Link",
				`<https://test-instance.com/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=1>; rel="prev",
				<https://test-instance.com/api/v3/enterprises/SGNL/secret-scanning/alerts?per_page=1&page=1>; rel="first"`,
			)
			w.Write([]byte(`[
				{
					"number": 1,
					"created_at": "2020-11-06T18:18:30Z",
					"url": "https://api.github.com/repos/owner/repo/secret-scanning/alerts/1",
					"html_url": "https://github.com/owner/repo/security/secret-scanning/1",
					"locations_url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/1/locations",
					"state": "open",
					"resolution": null,
					"resolved_at": null,
					"resolved_by": null,
					"secret_type": "mailchimp_api_key",
					"secret_type_display_name": "Mailchimp API Key",
					"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2",
					"repository": {
					  "id": 1296269,
					  "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
					  "name": "Hello-World",
					  "full_name": "octocat/Hello-World",
					  "owner": {
						"login": "octocat",
						"id": 1,
						"node_id": "MDQ6VXNlcjE=",
						"avatar_url": "https://github.com/images/error/octocat_happy.gif",
						"gravatar_id": "",
						"url": "https://api.github.com/users/octocat",
						"html_url": "https://github.com/octocat",
						"followers_url": "https://api.github.com/users/octocat/followers",
						"following_url": "https://api.github.com/users/octocat/following{/other_user}",
						"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
						"organizations_url": "https://api.github.com/users/octocat/orgs",
						"repos_url": "https://api.github.com/users/octocat/repos",
						"events_url": "https://api.github.com/users/octocat/events{/privacy}",
						"received_events_url": "https://api.github.com/users/octocat/received_events",
						"type": "User",
						"site_admin": false
					  },
					  "private": false,
					  "html_url": "https://github.com/octocat/Hello-World",
					  "description": "This your first repo!",
					  "fork": false,
					  "url": "https://api.github.com/repos/octocat/Hello-World",
					  "archive_url": "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
					  "assignees_url": "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
					  "blobs_url": "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
					  "branches_url": "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
					  "collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
					  "comments_url": "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
					  "commits_url": "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
					  "compare_url": "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
					  "contents_url": "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
					  "contributors_url": "https://api.github.com/repos/octocat/Hello-World/contributors",
					  "deployments_url": "https://api.github.com/repos/octocat/Hello-World/deployments",
					  "downloads_url": "https://api.github.com/repos/octocat/Hello-World/downloads",
					  "events_url": "https://api.github.com/repos/octocat/Hello-World/events",
					  "forks_url": "https://api.github.com/repos/octocat/Hello-World/forks",
					  "git_commits_url": "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
					  "git_refs_url": "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
					  "git_tags_url": "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
					  "issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
					  "issue_events_url": "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
					  "issues_url": "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
					  "keys_url": "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
					  "labels_url": "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
					  "languages_url": "https://api.github.com/repos/octocat/Hello-World/languages",
					  "merges_url": "https://api.github.com/repos/octocat/Hello-World/merges",
					  "milestones_url": "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
					  "notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
					  "pulls_url": "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
					  "releases_url": "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
					  "stargazers_url": "https://api.github.com/repos/octocat/Hello-World/stargazers",
					  "statuses_url": "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
					  "subscribers_url": "https://api.github.com/repos/octocat/Hello-World/subscribers",
					  "subscription_url": "https://api.github.com/repos/octocat/Hello-World/subscription",
					  "tags_url": "https://api.github.com/repos/octocat/Hello-World/tags",
					  "teams_url": "https://api.github.com/repos/octocat/Hello-World/teams",
					  "trees_url": "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
					  "hooks_url": "https://api.github.com/repos/octocat/Hello-World/hooks"
					},
					"push_protection_bypassed_by": null,
					"push_protection_bypassed": false,
					"push_protection_bypassed_at": null,
					"resolution_comment": null,
					"validity": "unknown"
				  }
			]`))

		// Secret Scanning Alerts only 1 page for org1.
		case "/api/v3/orgs/org1/secret-scanning/alerts?per_page=1":
			w.Header().Add("Link",
				`<https://test-instance.com/api/v3/orgs/org1/secret-scanning/alerts?per_page=1&page=2>; rel="last"`,
			)
			w.Write([]byte(`[
				{
					"number": 2,
					"created_at": "2020-11-06T18:48:51Z",
					"url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2",
					"html_url": "https://github.com/owner/private-repo/security/secret-scanning/2",
					"locations_url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/2/locations",
					"state": "resolved",
					"resolution": "false_positive",
					"resolved_at": "2020-11-07T02:47:13Z",
					"resolved_by": {
					  "login": "monalisa",
					  "id": 2,
					  "node_id": "MDQ6VXNlcjI=",
					  "avatar_url": "https://alambic.github.com/avatars/u/2?",
					  "gravatar_id": "",
					  "url": "https://api.github.com/users/monalisa",
					  "html_url": "https://github.com/monalisa",
					  "followers_url": "https://api.github.com/users/monalisa/followers",
					  "following_url": "https://api.github.com/users/monalisa/following{/other_user}",
					  "gists_url": "https://api.github.com/users/monalisa/gists{/gist_id}",
					  "starred_url": "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
					  "subscriptions_url": "https://api.github.com/users/monalisa/subscriptions",
					  "organizations_url": "https://api.github.com/users/monalisa/orgs",
					  "repos_url": "https://api.github.com/users/monalisa/repos",
					  "events_url": "https://api.github.com/users/monalisa/events{/privacy}",
					  "received_events_url": "https://api.github.com/users/monalisa/received_events",
					  "type": "User",
					  "site_admin": true
					},
					"secret_type": "adafruit_io_key",
					"secret_type_display_name": "Adafruit IO Key",
					"secret": "aio_XXXXXXXXXXXXXXXXXXXXXXXXXXXX",
					"repository": {
					  "id": 1296269,
					  "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
					  "name": "Hello-World",
					  "full_name": "octocat/Hello-World",
					  "owner": {
						"login": "octocat",
						"id": 1,
						"node_id": "MDQ6VXNlcjE=",
						"avatar_url": "https://github.com/images/error/octocat_happy.gif",
						"gravatar_id": "",
						"url": "https://api.github.com/users/octocat",
						"html_url": "https://github.com/octocat",
						"followers_url": "https://api.github.com/users/octocat/followers",
						"following_url": "https://api.github.com/users/octocat/following{/other_user}",
						"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
						"organizations_url": "https://api.github.com/users/octocat/orgs",
						"repos_url": "https://api.github.com/users/octocat/repos",
						"events_url": "https://api.github.com/users/octocat/events{/privacy}",
						"received_events_url": "https://api.github.com/users/octocat/received_events",
						"type": "User",
						"site_admin": false
					  },
					  "private": false,
					  "html_url": "https://github.com/octocat/Hello-World",
					  "description": "This your first repo!",
					  "fork": false,
					  "url": "https://api.github.com/repos/octocat/Hello-World",
					  "archive_url": "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
					  "assignees_url": "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
					  "blobs_url": "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
					  "branches_url": "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
					  "collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
					  "comments_url": "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
					  "commits_url": "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
					  "compare_url": "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
					  "contents_url": "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
					  "contributors_url": "https://api.github.com/repos/octocat/Hello-World/contributors",
					  "deployments_url": "https://api.github.com/repos/octocat/Hello-World/deployments",
					  "downloads_url": "https://api.github.com/repos/octocat/Hello-World/downloads",
					  "events_url": "https://api.github.com/repos/octocat/Hello-World/events",
					  "forks_url": "https://api.github.com/repos/octocat/Hello-World/forks",
					  "git_commits_url": "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
					  "git_refs_url": "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
					  "git_tags_url": "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
					  "issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
					  "issue_events_url": "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
					  "issues_url": "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
					  "keys_url": "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
					  "labels_url": "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
					  "languages_url": "https://api.github.com/repos/octocat/Hello-World/languages",
					  "merges_url": "https://api.github.com/repos/octocat/Hello-World/merges",
					  "milestones_url": "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
					  "notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
					  "pulls_url": "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
					  "releases_url": "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
					  "stargazers_url": "https://api.github.com/repos/octocat/Hello-World/stargazers",
					  "statuses_url": "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
					  "subscribers_url": "https://api.github.com/repos/octocat/Hello-World/subscribers",
					  "subscription_url": "https://api.github.com/repos/octocat/Hello-World/subscription",
					  "tags_url": "https://api.github.com/repos/octocat/Hello-World/tags",
					  "teams_url": "https://api.github.com/repos/octocat/Hello-World/teams",
					  "trees_url": "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
					  "hooks_url": "https://api.github.com/repos/octocat/Hello-World/hooks"
					},
					"push_protection_bypassed_by": {
					  "login": "monalisa",
					  "id": 2,
					  "node_id": "MDQ6VXNlcjI=",
					  "avatar_url": "https://alambic.github.com/avatars/u/2?",
					  "gravatar_id": "",
					  "url": "https://api.github.com/users/monalisa",
					  "html_url": "https://github.com/monalisa",
					  "followers_url": "https://api.github.com/users/monalisa/followers",
					  "following_url": "https://api.github.com/users/monalisa/following{/other_user}",
					  "gists_url": "https://api.github.com/users/monalisa/gists{/gist_id}",
					  "starred_url": "https://api.github.com/users/monalisa/starred{/owner}{/repo}",
					  "subscriptions_url": "https://api.github.com/users/monalisa/subscriptions",
					  "organizations_url": "https://api.github.com/users/monalisa/orgs",
					  "repos_url": "https://api.github.com/users/monalisa/repos",
					  "events_url": "https://api.github.com/users/monalisa/events{/privacy}",
					  "received_events_url": "https://api.github.com/users/monalisa/received_events",
					  "type": "User",
					  "site_admin": true
					},
					"push_protection_bypassed": true,
					"push_protection_bypassed_at": "2020-11-06T21:48:51Z",
					"resolution_comment": "Example comment",
					"validity": "active"
				  }
			]`))

		case "/api/v3/orgs/org2/secret-scanning/alerts?per_page=1":
			w.Header().Add("Link",
				`<https://test-instance.com/api/v3/orgs/org2/secret-scanning/alerts?per_page=1&page=2>; rel="last"`,
			)
			w.Write([]byte(`[
				{
					"number": 1,
					"created_at": "2020-11-06T18:18:30Z",
					"url": "https://api.github.com/repos/owner/repo/secret-scanning/alerts/1",
					"html_url": "https://github.com/owner/repo/security/secret-scanning/1",
					"locations_url": "https://api.github.com/repos/owner/private-repo/secret-scanning/alerts/1/locations",
					"state": "open",
					"resolution": null,
					"resolved_at": null,
					"resolved_by": null,
					"secret_type": "mailchimp_api_key",
					"secret_type_display_name": "Mailchimp API Key",
					"secret": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX-us2",
					"repository": {
					  "id": 1296269,
					  "node_id": "MDEwOlJlcG9zaXRvcnkxMjk2MjY5",
					  "name": "Hello-World",
					  "full_name": "octocat/Hello-World",
					  "owner": {
						"login": "octocat",
						"id": 1,
						"node_id": "MDQ6VXNlcjE=",
						"avatar_url": "https://github.com/images/error/octocat_happy.gif",
						"gravatar_id": "",
						"url": "https://api.github.com/users/octocat",
						"html_url": "https://github.com/octocat",
						"followers_url": "https://api.github.com/users/octocat/followers",
						"following_url": "https://api.github.com/users/octocat/following{/other_user}",
						"gists_url": "https://api.github.com/users/octocat/gists{/gist_id}",
						"starred_url": "https://api.github.com/users/octocat/starred{/owner}{/repo}",
						"subscriptions_url": "https://api.github.com/users/octocat/subscriptions",
						"organizations_url": "https://api.github.com/users/octocat/orgs",
						"repos_url": "https://api.github.com/users/octocat/repos",
						"events_url": "https://api.github.com/users/octocat/events{/privacy}",
						"received_events_url": "https://api.github.com/users/octocat/received_events",
						"type": "User",
						"site_admin": false
					  },
					  "private": false,
					  "html_url": "https://github.com/octocat/Hello-World",
					  "description": "This your first repo!",
					  "fork": false,
					  "url": "https://api.github.com/repos/octocat/Hello-World",
					  "archive_url": "https://api.github.com/repos/octocat/Hello-World/{archive_format}{/ref}",
					  "assignees_url": "https://api.github.com/repos/octocat/Hello-World/assignees{/user}",
					  "blobs_url": "https://api.github.com/repos/octocat/Hello-World/git/blobs{/sha}",
					  "branches_url": "https://api.github.com/repos/octocat/Hello-World/branches{/branch}",
					  "collaborators_url": "https://api.github.com/repos/octocat/Hello-World/collaborators{/collaborator}",
					  "comments_url": "https://api.github.com/repos/octocat/Hello-World/comments{/number}",
					  "commits_url": "https://api.github.com/repos/octocat/Hello-World/commits{/sha}",
					  "compare_url": "https://api.github.com/repos/octocat/Hello-World/compare/{base}...{head}",
					  "contents_url": "https://api.github.com/repos/octocat/Hello-World/contents/{+path}",
					  "contributors_url": "https://api.github.com/repos/octocat/Hello-World/contributors",
					  "deployments_url": "https://api.github.com/repos/octocat/Hello-World/deployments",
					  "downloads_url": "https://api.github.com/repos/octocat/Hello-World/downloads",
					  "events_url": "https://api.github.com/repos/octocat/Hello-World/events",
					  "forks_url": "https://api.github.com/repos/octocat/Hello-World/forks",
					  "git_commits_url": "https://api.github.com/repos/octocat/Hello-World/git/commits{/sha}",
					  "git_refs_url": "https://api.github.com/repos/octocat/Hello-World/git/refs{/sha}",
					  "git_tags_url": "https://api.github.com/repos/octocat/Hello-World/git/tags{/sha}",
					  "issue_comment_url": "https://api.github.com/repos/octocat/Hello-World/issues/comments{/number}",
					  "issue_events_url": "https://api.github.com/repos/octocat/Hello-World/issues/events{/number}",
					  "issues_url": "https://api.github.com/repos/octocat/Hello-World/issues{/number}",
					  "keys_url": "https://api.github.com/repos/octocat/Hello-World/keys{/key_id}",
					  "labels_url": "https://api.github.com/repos/octocat/Hello-World/labels{/name}",
					  "languages_url": "https://api.github.com/repos/octocat/Hello-World/languages",
					  "merges_url": "https://api.github.com/repos/octocat/Hello-World/merges",
					  "milestones_url": "https://api.github.com/repos/octocat/Hello-World/milestones{/number}",
					  "notifications_url": "https://api.github.com/repos/octocat/Hello-World/notifications{?since,all,participating}",
					  "pulls_url": "https://api.github.com/repos/octocat/Hello-World/pulls{/number}",
					  "releases_url": "https://api.github.com/repos/octocat/Hello-World/releases{/id}",
					  "stargazers_url": "https://api.github.com/repos/octocat/Hello-World/stargazers",
					  "statuses_url": "https://api.github.com/repos/octocat/Hello-World/statuses/{sha}",
					  "subscribers_url": "https://api.github.com/repos/octocat/Hello-World/subscribers",
					  "subscription_url": "https://api.github.com/repos/octocat/Hello-World/subscription",
					  "tags_url": "https://api.github.com/repos/octocat/Hello-World/tags",
					  "teams_url": "https://api.github.com/repos/octocat/Hello-World/teams",
					  "trees_url": "https://api.github.com/repos/octocat/Hello-World/git/trees{/sha}",
					  "hooks_url": "https://api.github.com/repos/octocat/Hello-World/hooks"
					},
					"push_protection_bypassed_by": null,
					"push_protection_bypassed": false,
					"push_protection_bypassed_at": null,
					"resolution_comment": null,
					"validity": "unknown"
				  }
			]`))
		}
	}
})
