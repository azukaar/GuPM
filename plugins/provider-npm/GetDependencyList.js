var result = [];

for(depName in PackageConfig.dependencies) {
    var depVersion = PackageConfig.dependencies[depName];
    result.push({
        provider: 'npm',
        name: depName,
        version: depVersion.replace(/(\d) ([\>\<\=\^\~\!])/g, '$1, $2')
    })
}

result;