# Dashboard

Without any doubt, this is the principal resource of Perses.

A `Dashboard` belongs to a `Project`. See the [project documentation](./project.md) to see how to create a project.

It is defined like that:

```yaml
kind: "Dashboard"
metadata:
  name: <string>
  project: <string>
spec: <dashboard_specification>
```

See the next section to get details about the `<dashboard_specification>`.

## Dashboard specification

```yaml
  # Metadata.name has some restrictions. For example, you can't use space there.
  # `display` allows to provide a rich name and a description for your dashboard.
  [ display: <display_spec> ]

  datasources:
    [ <string>: <datasource_spec> ]

  # `variables` is the list of dashboard variables. A variable can be referenced by the different panels and/or by other variables.
  [ variables: <variable_spec> ]

  # `panels` is a map where the key is the reference of the panel. The value is the actual panel definition that describes
  # the kind of chart this panel is using. A panel can only hold one chart.
  panels:
    [ <string>: <panel_spec> ]

  # `layouts` is the list of layouts. A layout describes how to display the list of panels. 
  # Indeed, in Perses the definition of a panel is uncorrelated from the definition of where to position it.
  layouts:
    - <layout_spec>

  # `duration` is the default time range to use on the initial load of the dashboard.
  [ duration: <duration> ]

  # `refreshInterval` is the default refresh interval to use on the initial load of the dashboard.
  [ refreshInterval: <duration> ]
```

A dashboard in its minimal definition only requires a panel and a layout.

### `<display_spec>`

This is the way to provide a rich name and a description for your dashboard. There is no restriction about the type of
characters you can use there.

```yaml
  # The new name of the dashboard. If set, it will replace `metadata.name` in the dashboard title in the UI.
  # Note that it cannot be used when you are querying the API. Only `metadata.name` can be used to reference the dashboard.
  # This is just for display purpose.
  [ name: <string> ]

  # The description of the dashboard.
  [ description: <string> ]
```

### `<datasource_spec>`

See the [datasource](./datasource.md) documentation.

### `<variable_spec>`

See the [variable](./variable.md) documentation.

### `<panel_spec>`

```yaml
kind: "Panel"
spec:
  display: <display_spec>

  # `plugin` is where you define the chart type to use.
  # The chart type chosen should match one of the chart plugins known to the Perses instance.
  plugin: <panel_plugin_spec>

  # `queries` is the list of queries to be executed by the panel. The available types of query are conditioned by the type of chart & the type of datasource used.
  queries:
    - [ <query_spec> ]
```

#### `<panel_plugin_spec>`

```yaml
  # `kind` is the plugin type of the panel. For example, `TimeSeriesChart`.
  kind: <string>

  # `spec` is the actual definition of the panel plugin. Each `kind` comes with its own `spec`.
  spec: <plugin_spec>
```

See the [panel](../plugin/panel.md) documentation to know more about the different panel plugins supported by Perses.

#### `<query_spec>`

```yaml
# kind` is the type of the query. For the moment we only support `TimeSeriesQuery`.
kind: <string>
spec:
  plugin: <query_plugin_spec>
```

##### `<query_plugin_spec>`

```yaml
  # `kind` is the plugin type matching the type of query. For example, `PrometheusTimeSeriesQuery` for the query type `TimeSeriesQuery`.
  kind: <string>

  # `spec` is the actual definition of the query. Each `kind` comes with its own `spec`.
  spec: <plugin_spec>
```

We are supporting only prometheus for the `TimeSeriesQuery` for the moment.
Please look at the [Prometheus plugin documentation](../plugin/prometheus.md#datasource) to know the spec for the `PrometheusTimeSeriesQuery`.

### `<layout_spec>`

```yaml
kind: "Grid"
spec:
  [ display: <grid_display_spec> ]
  items:
    [ - <grid_item_spec> ]
```

### `<grid_item_spec>`

```yaml
x: <int>
y: <int>
width: <int>
height: <int>
content:
  "$ref": <json_panel_ref>
```

Example:

```yaml
kind: "Grid"
spec:
  display:
    title: "Row 1"
    collapse:
      open: true
  items:
    - x: 0
      y: 0
      width: 2
      height: 3
      content:
        "$ref": "#/spec/panels/statRAM"
    - x: 0
      y: 4
      width: 2
      height: 3
      content:
        $ref": "#/spec/panels/statTotalRAM"
```

## API definition

### Get a list of `Dashboard`

```bash
GET /api/v1/projects/<project_name>/dasbhoards
```

URL query parameters:

- name = `<string>` : filters the list of dashboards based on their name (prefix match).

### Get a single `Dashboard`

```bash
GET /api/v1/projects/<project_name>/dasbhoards/<dasbhoard_name>
```

### Create a single `Dashboard`

```bash
POST /api/v1/projects/<project_name>/dashboards
```

### Update a single `Dashboard`

```bash
PUT /api/v1/projects/<project_name>/dasbhoards/<dasbhoard_name>
```

### Delete a single `Dashboard`

```bash
DELETE /api/v1/projects/<project_name>/dasbhoards/<dasbhoard_name>
```
