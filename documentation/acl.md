# Role Based Access Control

Every role has multiple statements

Statements consist of what the user has access to

#### Statements

Statements define what an entity will have access to but they are not tied to a specific entity.

Every statement will contain the following data:

- `ResourceName`
- `Action`

##### ResourceName

The resource name denotes the name of a specific resource that this statement applies to.

For example, if there should be some access given or taken away from an organization. The resource name would denote `organization`.

##### ResourceID

The resource ID denotes the actual ID for the named resource that the statement applies to.

A valid resource ID is either an asterisk, or a uuid. The asterisk (`*`) denotes that given actions in the statement apply to all of the resources with the given name.

##### Action

Actions is an array of values that denote the actions that can be applied to a given resource. For example, if there should be read and write actions allowed for this resource, the actions would consist of `read` and `write`.

##### Other Quirks

- There are two actions that are implicitly allowed by anyone. That is creating a new user and creating a new organization. 
- All resources have statements associated with them. For example every single document has a 


##### Examples

Read and write access for an organization
```
{
    ResourceName: "organization",
    Actions: "read"
}
```

Ability to create a folder in an organization
```
{
    ResourceName: "organization"
    ResourceID: "12345"
    Actions: "create:folder"
}
```

Ability to create a document in a folder
```
{
    ResourceName: "folder"
    ResourceId: "12345"
    Actions: "create:document"
}
```

#### Policy

A policy applies to a specific resource. A policy is a collection of statements.


User A can read documents A, B, C in folder A, but user B can not only see D

Statements are defined for every resource that is created. A policy then consists of a collection of statements. Policies are then assigned to individual users.

This gives granular control (in the future) over individual files.

For example, if you want to grant a group of users access to all documents in a folder. You would have a policy `Folder 5 Access`, that grants read access to the folder and to each document. The policy would contain read statements for all documents and a read statement for the folder (example below). If you wanted to remoove read access to a document for all users. You wouldd simply find the statement that grants read access to the document, and remove it from all policies. This would automatically remove read access for the document for everyone in that existing policy. A new policy can be created for just the document as well. And this can be assigned to another user. This gives heavy customization to who can access what in what ways.

Templates can also be created that denote a "role". These templates simply consist of the statements (and policy) that should be created for a given resource to denote a specific level of access. E.g. if a user is granted "contributor" access to a document. They would have a policy created that gives `read` and `write` access. For now these templates are defined in code.

```
{
    UserId: "12345",
    Policies: [
        {
            Name: "Folder 5 Access"
            Statements: [
                {
                    ResourceName: "folder"
                    ResourceId: "12345"
                    Action: "read" 
                },
                {
                    ResourceName: "document"
                    ResourceId: "12345"
                    Action: "read" 
                },
                {
                    ResourceName: "document"
                    ResourceId: "54321"
                    Action: "read" 
                },
            ]
        }
    ]

}
```

#### Tables

##### Statement

- id
- resourceName
- resourceId
- action

##### Policy

- id
- name

##### PolicyStatement

- policy_id
- statement_id

##### UserPolicy

- user_id
- policy_id

##### User

...the normal user values