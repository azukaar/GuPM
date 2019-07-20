var name = Dependency.name;

if(name.split(':').length > 1) {
    var version = Dependency.version
    var names = name.split(':')

    // exact version
    if (version.match(/^\d+\.\d+\.\d+/) && !version.match(/\sx/)) {
    }
    
    // test ranges
    else {
        var repoList = 'https://' + names[0] + '/gupm_repo.json'
        var payload = httpGetJson(repoList);

        if(payload.packages[names[1]] && payload.packages[names[1]].length) {
            var versionList = payload.packages[names[1]];
            version = semverLatestInRange(version, versionList);
        } else {
            console.error("Package "+names[1]+" not found in " + names[0])
            exit(1)
        }
    }



    var realName = names[1]
    var namespace = ''
    if( realName.split('/').length > 1) {
        namespace = (realName.split('/')[0] + '/').replace("OS", _OSNAME)
        realName = realName.split('/')[1]
    }
    Dependency.url = 'https://' + names[0] + '/' + namespace + realName + '/'  + version + '/' + realName + '-'  + version + '.tgz'

    Dependency.version = version
} else {
    Dependency.url = 'https://' + name
    Dependency.name = Dependency.name.replace(/\//g, "-")
}

Dependency;
