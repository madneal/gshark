$('input[type=radio][name=status]').change(function() {
	let status;
	if (this.value === 'unhandled') {
	    status = 0;
    } else if (this.value === 'confirmed') {
	    status = 1;
    } else {
	    status = 2;
    }
    window.location.href = ('/admin/reports/github/query/' + status);
})