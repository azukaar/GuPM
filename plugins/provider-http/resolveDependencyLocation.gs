var name = Dependency.name;

if(name.split(':').length > 1) {
    var version = Dependency.version
    var names = name.split(':')

    // exact version
    if (version.match(/^\d+\.\d+\.\d+/) && !version.match(/\sx/)) {
    }
    
    // test ranges
    else {
        var repoList = 'http://' + names[0] + '/gupm_repo.json'
        var payload = httpGetJson(repoList);
        var versionList = payload.packages[names[1]];
        version = semverLatestInRange(version, versionList);
    }


    Dependency.url = 'http://' + names[0] + '/' + names[1] + '/'  + version + '/' + names[1] + '-'  + version + '.tgz'
    Dependency.name = names[1]
    Dependency.version = version
} else {
    Dependency.url = 'http://' + name
    Dependency.name = Dependency.name.replace(/\//g, "-")
}

Dependency;
