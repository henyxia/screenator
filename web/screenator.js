/* helpers */
function closeAllSectionExcept(secName) {
    var sections = document.getElementsByTagName("section");
    for(var i = 0; i < sections.length; i++){
        if (sections[i].id != secName)
            sections[i].style.display = 'none'
        else
            sections[i].style.display = 'block'
    }

}

function httpGetAsync(theUrl, callback)
{
    var xmlHttp = new XMLHttpRequest();
    xmlHttp.onreadystatechange = function() {
    if (xmlHttp.readyState == 4 && xmlHttp.status == 200)
        callback(xmlHttp.responseText);
    }
    xmlHttp.open("GET", theUrl, true); // true for asynchronous
    xmlHttp.send(null);
}

function apiGet(item, cb) {
    httpGetAsync(urlListToUrl(item), cb);
}

function urlListToUrl(item) {
    var urls = {
        site: '/site',
    }
    return urls[item];
}

/* router like functions */
function goHome() {
    closeAllSectionExcept('home');
}

function goSiteList() {
    closeAllSectionExcept('site');
    apiGet('site', function(content) {
        sites = JSON.parse(content)
        tbody = document.getElementById('site_list');
        newtbody = ''
        sites.forEach(function (site) {
            item = "<tr><td>";
            item += site.ID;
            item += "</td><td>";
            item += site.Name;
            item += "</td></tr>";
            newtbody += item;
            console.debug(site);
        })
        tbody.innerHTML = newtbody;
    })
}

function goLocationList() {
    closeAllSectionExcept('location');
}

function goDeviceList() {
    closeAllSectionExcept('device');
}

function goContentList() {
    closeAllSectionExcept('content');
}

function goDisplayList() {
    closeAllSectionExcept('display');
}
