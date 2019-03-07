$('input[type=radio][name=status]').change(function() {
	let status;
	if (this.value === 'unhandled') {
	    status = 0;
    } else if (this.value === 'confirmed') {
	    status = 1;
    } else {
	    status = 2;
    }
    window.location.href = window.location.origin + '/admin/reports/github/query/' + status
    // if (window.location.pathname.indexOf("query") == -1) {
	 //    window.location.href = window.location.href + "query/" + status;
    // } else {
    //     window.location.href = window.location.origin + window.location.pathname.replace(/\d/, status);
    // }
})
var url = window.location;
// for sidebar menu but not for treeview submenu
$('ul.sidebar-menu a').filter(function() {
    return this.href == url;
}).parent().siblings().removeClass('active').end().addClass('active');
// for treeview which is like a submenu
$('ul.treeview-menu a').filter(function() {
    return this.href == url;
}).parentsUntil(".sidebar-menu > .treeview-menu").siblings().removeClass('active').end().addClass('active');