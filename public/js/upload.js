document.querySelector('.file input').addEventListener('change', function () {
    var blobSlice = File.prototype.slice || File.prototype.mozSlice || File.prototype.webkitSlice,
        file = this.files[0],
        chunkSize = 2097152,                             // Read in chunks of 2MB
        chunks = Math.ceil(file.size / chunkSize),
        currentChunk = 0,
        spark = new SparkMD5.ArrayBuffer(),
        fileReader = new FileReader();

    fileReader.onload = function (e) {
        console.log('read chunk nr', currentChunk + 1, 'of', chunks);
        spark.append(e.target.result);                   // Append array buffer
        currentChunk++;
        let hash;

        if (currentChunk < chunks) {
            loadNext();
        } else {
            hash = spark.end();
            // alert(hash);
            $.get("/admin/app/detect", {"hash": hash}, function (data) {
                if (data.isExist) {
                    alert("该 APP 存在")
                    window.location.href = "/admin/app/appid/?id=" + data.id;
                } else {
                    alert("该 APP 不存在");
                    window.location.href = "admin/app/new";
                }
            })
        }
    };

    fileReader.onerror = function () {
        console.warn('oops, something went wrong.');
    };

    function loadNext() {
        var start = currentChunk * chunkSize,
            end = ((start + chunkSize) >= file.size) ? file.size : start + chunkSize;

        fileReader.readAsArrayBuffer(blobSlice.call(file, start, end));
    }

    loadNext();
});
