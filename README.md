# Helm tool

Tool for generating documentation, generating values.schema.json and linting values.yaml. It works by reading the yaml content and the surrounding comments to infer information.

## Usage

- `helm-tool schema -i values.yaml > values.schema.json` - Generate a values.schema.json file based on the properties in values.yaml
- `helm-tool lint -i values.yaml -d templates -e values.linter.exceptions` - Lint the values.yaml properties based on what properties are used in the template (imperfect linter, might miss errors or report false positives)

There are two commands that can be used to generate documentation, `helm-tool render` and `helm-tool inject`.

- `helm-tool render` - The render command will simply render the markdown to the stdout
- `helm-tool inject` - The inject command will inject the generated documentation into an existing markdown file, it will look for the `## Properties` header and inject the documentation between it and the next header. This can be useful for keeping a chart README up to date.

## Customising the output

### Sections

Documentation can be divided up into sections through the `+docs:section` tag, for example:

```yaml
# +docs:section=Global
# This contains all parameters that are used by all deployments in this chart.

# Foo parameter description
foo: bar

# +docs:section=Application
# Application specific parameters

# Baz parameter description
baz: qux
```

Would produce the following markdown:

```markdown
### Global
This contains all parameters that are used by all deployments in this chart.

|property|description|type|default|
|--|--|--|--|
|`foo`|<p>Foo parameter description</p>|`string`|<pre>bar</pre>|

### Application
Application specific parameters

|property|description|type|default|
|--|--|--|--|
|`baz`|<p>Baz parameter description</p>|`string`|<pre>qux</pre>|
```

### Undefaulted properties

Often helm values files have properties that do not require a default value commented out, this tool can find those 
by marking them with the `+docs:property` tag. 

For example:

```yaml
# +docs:property
# Property description here
# foo: bar
```

Would produce the following markdown:

```markdown
|property|description|type|default|
|--|--|--|--|
|`foo`|<p>Property description here</p>|`string`|<pre>undefined</pre>|
```


The detected name and type is not always correct, these can be provided using tags. For example:
```yaml
# +docs:property=foo
# +docs:type=string
# Property description here
```

```markdown
|property|description|type|default|
|--|--|--|--|
|`foo`|<p>Property description here</p>|`string`|<pre>undefined</pre>|
```

### Tags

Tags are used to alter how the documentation is generated. They are comments that exist within a comment block

- `+docs:section=<name>` - Creates a new documentation section
- `+docs:property` - Marks the field as a property that needs documentation
- `+docs:ignore` - Ignore the field, not generating documentation, not used for linting or json schema generation
- `+docs:hidden` - Hide the field from the documentation, but still use it for linting and json schema generation
- `+docs:type=<type>` - Override the type information for the property
- `+docs:default=<default>` - Override the default value for the property
