$(document).ready(function() {
    $('tr', '#tracktable').click(function(){
        $('input[name="track"]').removeAttr('checked')
        $(this).closest('tr').find('input[name="track"]').prop('checked', true);
    })

    $('#track-play-button').click(function(){
        ControlTrack();
    })
});


function ControlTrack() {
	try {
		$('#stations-view').html('Try to execute');
		url = "/track/control/" + $("input[name='track']:checked").val();
		var postData = {};
		var service = new ServiceResult();
        service.getJSONDataRaw(url, postData, ControlTrack_Callback);
    }
    catch (e) {
        alert(e);
    }
}

function ControlTrack_Callback() {
	try {
        if (this.Data!=null && this.Data.Result==1) {
            $('#stations-view-json').html('Turn On counter');
        } else if (this.Data!=null && this.Data.Result==2) {
            $('#stations-view-json').html('Turn Off counter');
        } else if (this.Data!=null && this.Data.Result==3) {
            $('#stations-view-json').html(JSON.stringify(this.Data));
        } else {
            $('#stations-view-json').html(JSON.stringify(this.Data));
        }
	}
	catch (e) {
        alert(e);
    }
}

function Standard_Callback() {
    try {
        alert(this.ResultString);
    }
    catch (e) {   
        alert(e);
    }
}

function Standard_ValidationCallback() {
    try {
        alert(this.ResultString);
    }
    catch (e) {   
        alert(e);
    }
}

function Standard_ErrorCallback() {
    try {
        alert(this.ResultString);
    }
    catch (e) {   
        alert(e);
    }
}
