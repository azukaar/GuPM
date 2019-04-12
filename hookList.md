# Provider

## InitProvider

Chances to grab local repo info (apt / brew / etc...)

```
input:
output: true
```

## getPackageDocumentation

In case the doc isn't available to FS, you can hook this and return the json doc

```
input: packageDocumentationPath
output: JSON Object
```

## postGetPackageDocumentation

If you need to transform the doc before getting anything

```
input: packageDocumentation
output: JSON Object
```

## getDependencyList

Hook this to return the json array of deps

```
input: packageDocumentation (transformed)
output: JSON Array of deps
```

## postGetDependencyList

 - parse semver range and get real versions
 - build dependency structure (ex NPM with clashing versions)
 
```
input: JSON Array of deps
output: JSON Array of deps
```

## getDepedency

Download the dependency using the downloadFile function

```
input: single dep
output: true
```

## postGetDependency

Chance to unpack / link / copy the dependency

```
input: path to dep in cache
output: path to installed dep
```

## postInstallation

Called after each dependency's installation + root package installation

```
input: path to installed dep
output: true
```

## finalHook

Called after everything

- chance to link binaries (npm)

# Destination

## InitDestination

## prepare 
Build / zip 

## prePublish

## publish

## postPublish

# Destination