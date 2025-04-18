# GitHub Adapter/SoR Documentation

## Overview

This document outlines the entity relationships and pagination sync flows for the Github adapter.

## Entity Structure

- Enterprise
  - Organizations
    - Users
    - Repositories
      - Collaborators (Users)
      - Labels
        - IssueLabels (Connection Entity for Labels <-> Issues)
        - PullRequestLabels (Connection Entity for Labels <-> PullRequests)
      - Issues
        - IssueParticipants (Connection Entity for Issue <-> User Participants)
        - IssueAssignees (Connection Entity for Issue <-> User Assignees)
      - PullRequests
        - PRCommits
        - PRReviews
        - PRAssignees (Connection Entity for PullRequest <-> User Assignees)
        - PRParticipants (Connection Entity for PullRequest <-> User Participants)
        - PRChangedFiles
    - Teams
      - TeamMembers (Child Connection Entity for Teams <-> TeamMembers: Users with a team role)
      - TeamRepositories (Child Connection Entity for Teams <-> TeamRepositories: Repositories that a team has permission for)

### Notes:

- **Repeat:** Repeated entity structures that have already been expanded.
- **Enterprise Slug:** This is required as a parameter for the enterprise GitHub GraphQL query which is used in every sync.
- **Organization Login:** Required for every sync of user-type entities to access the 'organizationVerifiedDomainEmails' attribute.
- **OrganizationUser Entity:** OrganizationUser is a 'member' entity that we use to build relationships between Organizations and Users. This entity is unique because of the 'organizationVerifiedDomainEmails' (OVDE) attribute. This attribute is how we create relationships between GitHub user entities to other SoRs. In order to access this attribute, we need to specify the 'login' parameter which takes an organization login. As a result, anytime we want to request this parameter, we must use two queries: The first is a query using the Enterprise 'slug' attribute to retrieve organizations. The second query is a query using the organization 'login' attribute to get users. In this second query, we will also use the 'login' attribute as the parameter for the OVDE attribute. See the Postman Collection for sample queries and examples.
- **OVDE Attribute Ingested as Child Entity:** The 'organizationVerifiedDomainEmails' (OVDE) attribute is how we create relationships between GitHub user entities to other SoRs. Since OVDE is a list of strings in the GitHub response, we want to create relationships to each of the verified emails. This attribute has extra post-processing to convert the list of strings into a list of json objects so it can be ingested as a child entity.
- **Child Entities:** TeamMembers and TeamRepositories are currently the only child entities. This is because they don't have requirements for pagination and provide the option to receive all associated members and repositories of a team in bulk. In addition, TeamMembers is a subset of OrgUsers so we will not need to request the 'OVDE' attribute when syncing Teams/TeamMembers. Instead, 'OVDE' will be populated during the Users sync.
- **Collaborators Entity:** Collaborators is also a user-type entity and has been declared as standalone. Collaborators is not a subset of Users because it can contain external collaborators that have been assigned to repositories. Traditionally, entities like Collaborators would have been declared as a child since it is a list of objects that are associated with Repositories. However, we've opted to sync it separately to give us the flexibility of receiving the 'OVDE' attribute for external collaborators in the future.
- **Ignored Entity Branches:** Certain branches of entities are ignored due to redundancy and limitations in accessing organizationVerifiedDomainEmails without a corresponding organization 'login' attribute. For example there is also an Enterprise.Users branch that is ignored.
- **Connection Entities** Entities such as OrganizationUser and RepositoryCollaborator are ingested to create relationships between their corresponding entities. These connection entities need to be created for entities that have many-to-many relationships. However, entities like Team can form relationships through the 'orgId' attribute that is added manually during ingestion since the same Team can not exist across multiple organizations.
- **Container vs. Collection** The concept of collections and members revolves around the necessity of multiple queries to retrieve member information. Initially, a query is made to obtain a single collection ID. Subsequently, another query is executed to fetch the members associated with that specific ID. In GitHub Adapter, OrganizationUser is a 'member' entity of the Organization 'collection' (see explanation above). Containers/Entries is syntax used during the parsing logic of GitHub adapter. Since the GitHub GraphQL response is made up of many nested layers, we've divided this logic semantically as Container -> Entry relationships. We expect each intermediate container to only have a single entry. For instance, in the context of a Repository, the Organization serves as a container, and each entry is a Repository object within that container.

## Pagination

Similar to other adapters, GitHub also has collection/member entities that require multiple queries to retrieve a single page during a sync. For example, Users must be split into a query to get organizations from an enterprise, and another query to get users from an organization (Collection and Member Query). This is necessary since user-type entities require an organization 'login' attribute to use as a parameter to fetch the organizationVerifiedDomainEmails for a user.

However, the GitHub adapter supports entities that are retrieved using both the GitHub GraphQL and REST endpoints. Pagination is handled differently for each case.

**GraphQL:**
The CompositeCursor.Cursor string now stores a base64 encoded CursorInfo struct which contains up to 3 layers of 'After' parameters to be used in pagination. These parameters are used to paginate through multiple layers of entities at the same time. In the future we may need to add more layers.

**REST:**
The CompositeCursor.Cursor string stores a link to the endpoint that should be used to retrieve the next page of entities.
