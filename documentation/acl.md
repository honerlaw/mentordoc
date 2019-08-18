# Role Based Access Control

Role based access control consists of permissions, roles, and then a mapping of roles to users.

### Permissions

```
Permission {
    ResourcePath string
    Action string
}
```

A permission just denotes the path to the resource that a given action applies to. Permissions are never assigned directly to a user, they are always attached to a role first, and then assigned to a user.

#### ResourcePath

The resource path can be just a single resource, or denote the hierarchical path for the resource. For example, the path could just be `organization`, or
it could be `organization:folder`, or `folder`. The resource id (talked about in user role mapping), will apply to the root of the resource path. For example,
if the path is `organization:folder`, the id would apply to `organization`.

#### Actions

- `view`
- `modify`
- `delete`
- `create:{child_resource_name}`

#### Examples

The following two permissions would denote that a user would be given access to view the organization and to create folders in that organization.

```
Permission {
    ResourcePath "organization"
    Action "view"
}

Permission {
    ResourcePath "organization"
    Action "create:folder"
}
```

### Roles

```
Role {
    Name string
}
```

Roles define a collection of permissions. Users are always assigned permissions through roles.

### User Role Mapping

```
UserRole {
    UserId string
    RoleId string
    ResourceId string
```

User role mapping simply maps a set of permissions (role) to a given user and resource. An example of this would be
mapping a role containing permissions throughout the organization to a given user.

#### User Id

This is simply the identifier for a specific user that the role will be applied to.

#### Role Id

The identifier for the set of permissions to grant to the user.

#### Resource Id

The identifier of the resource that the role will be applied to.

### Examples

Organization Owner of Organization 54321
```
Permission {
    ResourcePath "organization"
    Action "view"
}

Permission {
    ResourcePath "organization"
    Action "modify"
}

Permission {
    ResourcePath "organization"
    Action "delete"
}

Permission {
    ResourcePath "organization"
    Action "create:folder"
}

Permission {
    ResourcePath "organization:folder"
    Action "view"
}

Permission {
    ResourcePath "organization:folder"
    Action "modify"
}

Permission {
    ResourcePath "organization:folder"
    Action "delete"
}

Permission {
    ResourcePath "organization:folder"
    Action "create:document"
}

Permission {
    ResourcePath "organization:folder:document"
    Action "view"
}

Permission {
    ResourcePath "organization:folder:document"
    Action "modify"
}

Permission {
    ResourcePath "organization:folder:document"
    Action "delete"
}

Role {
    Id "50"
    Name "organization:owner"
}

UserRole {
    UserId "12345"
    RoleId "50"
    ResourceId "54321"
}
```

Document Viewer of Document 54321
```
Permission {
    ResourcePath "document"
    Action "view"
}

Role {
    Id "50"
    Name "document:viewer"
}

UserRole {
    UserId "12345"
    RoleId "50"
    ResourceId "54321"
}
```

Folder Viewer and Document Owner of Folder 54321
```
Permission {
    ResourcePath "folder"
    Action "view"
}

Permission {
    ResourcePath "folder:document"
    Action "view"
}

Permission {
    ResourcePath "folder:document"
    Action "modify"
}

Permission {
    ResourcePath "folder:document"
    Action "delete"
}

Role {
    Id "50"
    Name "folder:document:owner"
}

UserRole {
    UserId "12345"
    RoleId "50"
    ResourceId "54321"
}
```