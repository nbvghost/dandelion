main.config(function ($routeProvider, $locationProvider, $provide, $httpProvider, $httpParamSerializerJQLikeProvider, $interpolateProvider) {

    $routeProvider.when("/content_item_list", {
        templateUrl: "content_templates/content_list_template",
        controller: "content_list_controller"
    });
    $routeProvider.when("/add_contents", {
        templateUrl: "content_templates/add_articles_template",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/edit_contents", {
        templateUrl: "content_templates/add_articles_template",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/add_blog", {
        templateUrl: "content_templates/add_blog_template",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/edit_blog", {
        templateUrl: "content_templates/add_blog_template",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/add_gallery", {
        templateUrl: "content_templates/add_gallery_template",
        controller: "content_add_gallery_controller"
    });
    $routeProvider.when("/edit_gallery", {
        templateUrl: "content_templates/add_articles_template",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/contents", {
        templateUrl: "content_templates/articles_template",
        controller: "content_articles_controller"
    });
    $routeProvider.when("/gallery", {
        templateUrl: "content_templates/articles_template",
        controller: "content_articles_controller"
    });
    $routeProvider.when("/blog", {
        templateUrl: "content_templates/articles_template",
        controller: "content_articles_controller"
    });
   /* $routeProvider.when("/content", {
        templateUrl: "content_templates/add_article_template",
        controller: "content_add_article_controller"
    });*/
    $routeProvider.when("/content_config", {
        templateUrl: "content_templates/content_config",
        controller: "content_config_controller"
    });


});


main.controller("content_config_controller", function ($http, $scope, $routeParams, $rootScope, $timeout, $location, Upload) {

    $scope.contentConfig = {}

    $scope.customerService = {
        SocialAccount:[]
    }


    $scope.socialAccount = {}


    $scope.tabIndex = 0
    $scope.selectTab = function (tabIndex) {
        $scope.tabIndex=tabIndex
    }
    $scope.deleteSocialAccount = function (index) {
        if(confirm("确定删除？")){
            $scope.customerService.SocialAccount.splice(index,1)
        }
    }
    $scope.editSocialAccount = function (index) {
        $scope.socialAccount = $scope.customerService.SocialAccount[index];
    }
    $scope.cancelSocialAccount = function () {
        $scope.socialAccount = {}
    }
    $scope.submitSocialAccount = function () {

        let Type=    $scope.socialAccount.Type||""
        let Account= $scope.socialAccount.Account||""
        if(Type.length>0 && Account.length>0) {
            let has =false
            for(let i=0;i<$scope.customerService.SocialAccount.length;i++){
                 let sa = $scope.customerService.SocialAccount[i]
                let accountInfo = sa.Type+sa.Account
                if(Type+Account===accountInfo){
                    has =true
                    break
                }
            }
            if(has===false){
                $scope.customerService.SocialAccount.push($scope.socialAccount)
                $scope.socialAccount={}
            }

        }
    }
    $scope.uploadCustomerServicePhoto = function (file, errFiles) {
        if (file) {
            const thumbnail = Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                const url = response.data.Path;
                $scope.customerService.Photo = url
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //upload_progress.progress('update progress',progress);
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        } else {
            if (errFiles.length > 0) {
                alert(JSON.stringify(errFiles));
            }

        }
    }
    $scope.submitCustomerService = function (){

        if($scope.customerService.SocialAccount.length===0){
            alert("请添加社交帐号")
            return
        }
        let Name =$scope.customerService.Name||""
        if(Name.length===0){
            alert("请填写名字")
            return
        }
        let Title =$scope.customerService.Title||""
        if(Title.length===0){
            alert("请填写头衔")
            return
        }

    }
    $scope.cancelCustomerService = function (){
        $("#CustomerServiceModal").modal({
            centered: false, closable: false, allowMultiple: false
        }).modal("hide");
    }
    $scope.showCustomerServiceModal = function (isAdd) {

        $("#CustomerServiceModal").modal({
            centered: false, closable: false, allowMultiple: false
        }).modal("show");
    }



    $timeout(function (){
        $('.tabular.menu .item').tab();
        $('.tabular.menu .item').tab('change tab', 'tab0');
        //$.tab('change tab', 'tab0');
    })

})
main.controller("content_articles_controller", function ($http, $scope, $routeParams, $rootScope, $timeout, $location, Upload) {
    $scope.ContentSubTypes = {};


    $scope.ContentSubTypeID;
    $scope.MContentSubTypeID;
    $scope.MContentSubTypeChildID;

    $scope.ContentItemID = $routeParams.ContentItemID;
    $scope.Type = $routeParams.Type;

    let table;

    $scope.listContentSubTypes = function () {
        //content/list
        $http.get("content/sub_type/list/all/" + $scope.ContentItemID).then(function (value) {

            $scope.ContentSubTypes = value.data.Data;

        });
    }
    $scope.listContentSubTypes();


    $scope.changeContentSubTypes = function () {
        $scope.ContentSubTypeID = $scope.MContentSubTypeID;
        table.ajax.reload();
    }
    $scope.changeContentSubTypeChilds = function () {
        $scope.ContentSubTypeID = $scope.MContentSubTypeChildID;
        table.ajax.reload();

    }

    $timeout(function () {

        table = $('#table_local').DataTable({
            "columns": [
                {data: "ID"},
                {data: "Title"},
                {data: "Author"},
                {data: "Look"},
                {data: "ContentItemID", visible: false},
                {data: "ContentSubTypeID", visible: false},
                {
                    data: null, className: "opera", orderable: false, render: function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>' +
                            '<button class="ui delete red mini button">删除</button>';

                    }
                }
            ],
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "order": [[1, "asc"]],
            "ajax": {
                "url": "content/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d) {
                    d.columns[4].search.value = parseInt($scope.ContentItemID).toString();


                    if ($scope.ContentSubTypeID) {
                        d.columns[5].search.value = parseInt($scope.ContentSubTypeID).toString();
                    }
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click', 'td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row(tr);
            console.log(row.data());
            window.location.href = "#!/edit_"+$scope.Type+"?ContentItemID=" + row.data().ContentItemID + "&ID=" + row.data().ID;
        });

        $('#table_local').on('click', 'td.opera .delete', function () {
                const tr = $(this).closest('tr');
                const row = table.row(tr);


                if (confirm("确定删除？")) {

                    var form = {};
                    form.ID = row.data().ID;
                    $http.post("content/delete", $.param(form), {
                        transformRequest: angular.identity,
                        headers: {"Content-Type": "application/x-www-form-urlencoded"}
                    }).then(function (data, status, headers, config) {

                        alert(data.data.Message);
                        table.ajax.reload();

                    });

                }
            }
        );
    })


});

main.controller("content_add_articles_controller", function ($http, $scope, $routeParams,$interval, $rootScope, $timeout, $location, Upload) {
    $scope.ContentSubTypes = [];

    $scope.ContentItemID = $routeParams.ContentItemID
    $scope.ContentSubTypeID = $routeParams.ContentSubTypeID
    $scope.ContentSubTypeChildID = $routeParams.ContentSubTypeChildID
    $scope.ArticleID = parseInt($routeParams.ID)//
    $scope.Single = $routeParams.Single//
    $scope.AutoSaveInfo = {Show:false,Time:"",Msg:""};

    let hasArticleID = true
    if ($scope.Single === "true") {
        hasArticleID = false
    }

    if (isNaN(parseInt($scope.ContentSubTypeID)) && isNaN(parseInt($scope.ContentSubTypeChildID))) {
        //alert("类别错误")
        //window.history.go(-1);
        //return
    }
    if (isNaN(parseInt($scope.ContentItemID))) {
        alert("类别错误")
        //window.history.go(-1);
        return
    }

    $scope.Article = {ContentItemID: $scope.ContentItemID};



    let articleContent = "";
    /*let runTime = $interval(async function (){

        $scope.AutoSaveInfo = {Show:false,Time:"",Msg:""};

        if(articleContent!==window.editor.getData()){

            let article =await $scope.saveArticle(true);

            articleContent=article.Content

            let queryParams=$location.search();
            queryParams.ID=article.ID;
            $location.search(queryParams);


            $scope.AutoSaveInfo = {Show:true,Time:article.UpdatedAt,Msg:""};

            $scope.ContentItemID = article.ContentItemID
            $scope.ArticleID = article.ID
            $scope.Single = false
            hasArticleID=true

            if(article.ContentSubTypeID===0){
                $scope.ContentSubTypeID = undefined
                $scope.ContentSubTypeChildID = undefined
            }else{

                await new Promise((resolve, reject) => {

                    $http.get("content/sub_type/get/" + article.ContentSubTypeID).then(async function (responea) {
                        let ContentSubType=responea.data.Data;
                        if(ContentSubType.ParentContentSubTypeID===0){
                            $scope.ContentSubTypeID = ContentSubType.ID
                            $scope.ContentSubTypeChildID =undefined
                        }else{
                            $scope.ContentSubTypeID = ContentSubType.ParentContentSubTypeID
                            $scope.ContentSubTypeChildID = ContentSubType.ID
                        }

                        resolve();

                    });
                });

                await loadArticle(hasArticleID);

            }
        }

    },1000000000);*/

    $scope.saveArticle = async function (isAutoSave) {

        /*$scope.Article.ContentItemID = parseInt($routeParams.ContentItemID);
        $scope.Article.Content = window.editor.getData();


        $scope.Article.ContentSubTypeID = parseInt($scope.Article.ContentSubTypeID)

        if (parseInt($scope.Article.ContentSubTypeChildID) > 0) {
            $scope.Article.ContentSubTypeID = parseInt($scope.Article.ContentSubTypeChildID)
        }

        if (!$scope.Article.ContentSubTypeID) {

            alert("请选择分类");
            return
        }*/

        return  new Promise((resolve, reject) => {
            let Article = angular.copy($scope.Article)

            Article.ContentItemID = parseInt($scope.ContentItemID);
            //Article.Content = window.editor.getData();
            Article.Content = quill.root.innerHTML;


            let ContentSubTypeID = 0

            if (isNaN(parseInt(Article.ContentSubTypeID)) === false) {
                ContentSubTypeID = parseInt(Article.ContentSubTypeID)
            }
            if (isNaN(parseInt(Article.ContentSubTypeChildID)) === false) {
                ContentSubTypeID = parseInt(Article.ContentSubTypeChildID)
            }


            Article.ContentSubTypeID = ContentSubTypeID

            $http.post("content/save", JSON.stringify(Article), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json"}
            }).then(function (data, status, headers, config) {
                console.log(data.data.Data);
                resolve(data.data.Data);
                if (data.data.Code === 0) {
                    if(!isAutoSave){
                        console.log(data);
                        alert(data.data.Message);
                        window.history.back();
                        window.close();
                    }

                }else{
                    alert(data.data.Message);
                }
            });
        });
    }


    $scope.UploadPictureImage = function (file, errFiles) {
        if (file) {
            const thumbnail = Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                const url = response.data.Path;
                $scope.Article.Picture = url;
            }, function (response) {
                console.log(response);
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //upload_progress.progress('update progress',progress);
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        } else {
            if (errFiles.length > 0) {
                alert(JSON.stringify(errFiles));
            }

        }
    }

    $scope.EditImages = [];
    $scope.UploadImages = function (progressID, files, errFiles) {

        const upload_progress = $(progressID);
        upload_progress.progress({
            duration: 100,
            total: 100,
            text: {
                active: '{value} of {total} done'
            }
        });

        upload_progress.progress('reset');
        //upload_progress.progress('update progress',50);

        if (files && files.length) {
            for (let i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up', data: {file: files[i]}}).then(function (response) {
                    const url = response.data.Path;

                    if ($scope.EditImages.indexOf(url) == -1) {
                        $scope.EditImages.push(url);
                    }

                }, function (response) {

                }, function (evt) {

                    const progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                    upload_progress.progress('update progress', progress);

                });
            }
        } else {
            UpImageError(errFiles);
        }
    }

    async function loadArticle(isLoadArticleID) {

        if (isLoadArticleID) {
            //通过Id来加载内容

            if (isNaN($scope.ArticleID)) {

                //如果没有传ID的话，说明是新加的内容
                return
            }

            await new Promise((resolve, reject) => {

                $http.get("content/multi/get/" + $scope.ArticleID).then(async function (responea) {

                    $scope.Article = responea.data.Data;

                    articleContent = $scope.Article.Content
                    //window.editor.setData($scope.Article.Content);
                    //quill.clipboard.dangerouslyPasteHTML(articleContent);
                    quill.root.innerHTML=articleContent

                    if ($scope.Article.ContentSubTypeID > 0) {

                        await new Promise((resolve1, reject1) => {
                            $http.get("content/sub_type/" + $scope.Article.ContentSubTypeID).then(function (responeb) {
                                const ContentSubType = responeb.data.Data.ContentSubType || {};
                                const ParentContentSubType = responeb.data.Data.ParentContentSubType || {};

                                $timeout(function () {


                                    if (parseInt(ContentSubType.ParentContentSubTypeID) > 0) {
                                        $scope.Article.ContentSubTypeID = ContentSubType.ParentContentSubTypeID + "";
                                        $scope.Article.ContentSubTypeChildID = ContentSubType.ID + "";
                                        // $scope.listContentSubTypeChilds($scope.Article.ContentItemID,$scope.MContentSubTypeID);
                                    } else {
                                        $scope.Article.ContentSubTypeID = ContentSubType.ID + "";
                                        $scope.Article.ContentSubTypeChildID = 0;
                                        // $scope.listContentSubTypeChilds($scope.Article.ContentItemID,$scope.MContentSubTypeID);
                                    }

                                    resolve1()

                                });


                            });
                        });
                    }


                    resolve()


                });

            });


        } else {
            await new Promise((resolve, reject) => {
                if ($scope.ContentItemID && ($scope.ContentSubTypeID || $scope.ContentSubTypeChildID)) {

                    let ContentSubTypeID = 0

                    if (isNaN(parseInt($scope.ContentSubTypeID)) === false) {
                        ContentSubTypeID = parseInt($scope.ContentSubTypeID)
                    }
                    if (isNaN(parseInt($scope.ContentSubTypeChildID)) === false) {
                        ContentSubTypeID = parseInt($scope.ContentSubTypeChildID)
                    }

                    $http.get("content/single/get/" + $scope.ContentItemID + "/" + ContentSubTypeID).then(function (responea) {

                        if (responea.data.Data.ID > 0) {
                            $scope.Article = responea.data.Data;

                            if ($scope.Article.ContentSubTypeID !== parseInt($scope.ContentSubTypeID) && $scope.Article.ContentSubTypeID !== parseInt($scope.ContentSubTypeChildID)) {

                                //alert("原内容类型与所选类型不匹配")
                                //return
                            }
                            //$scope.ContentSubTypeID = $scope.Article.ContentSubTypeID
                            //$scope.ContentSubTypeChildID = $scope.Article.ContentSubTypeChildID

                            $scope.Article.ContentSubTypeID = $scope.ContentSubTypeID || 0
                            $scope.Article.ContentSubTypeChildID = $scope.ContentSubTypeChildID || 0

                            //quill.clipboard.dangerouslyPasteHTML($scope.Article.Content);



                            articleContent = $scope.Article.Content
                            //quill.clipboard.dangerouslyPasteHTML(articleContent);
                            quill.root.innerHTML=articleContent

                            resolve()
                        } else {
                            resolve()
                        }

                    });
                } else {
                    resolve()
                }
            });
        }


    }

    let quill;

    $timeout(async function () {

        quill = new Quill('#editor-container',{
            modules: {
                formula: true,
                syntax: true,
                toolbar: '#toolbar-container',
                imageResize:{},
                //markdownShortcuts:{},
                imageUploader: {
                    upload: file => {
                        return new Promise((resolve, reject) => {
                            Upload.upload({url: '/file/up', data: {file: file}}).then(function (response) {
                                const url = response.data.Path;

                                resolve( '/file/load?path='+url)

                            }, function (evt) {

                                reject(evt)

                            });
                        });
                    }
                }
            },
            handlers: {

            },
            placeholder: 'Compose an epic...',
            theme: 'snow'
        });

        await new Promise((resolve, reject) => {
            $http.get("content/sub_type/list/all/" + $routeParams.ContentItemID).then(function (value) {
                $scope.ContentSubTypes = value.data.Data || [];

                $scope.Article.ContentSubTypeID = $scope.ContentSubTypeID
                $scope.Article.ContentSubTypeChildID = $scope.ContentSubTypeChildID
                console.log($scope.Article.ContentSubTypeID, $scope.Article.ContentSubTypeChildID)
                resolve()

            });
        })


        class MyUploadAdapter {
            constructor(loader) {
                // The file loader instance to use during the upload.
                this.loader = loader;
            }

            // Starts the upload process.
            upload() {
                return this.loader.file
                    .then(file => new Promise((resolve, reject) => {
                        this._initRequest();
                        this._initListeners(resolve, reject, file);
                        this._sendRequest(file);
                    }));
            }

            // Aborts the upload process.
            abort() {
                if (this.xhr) {
                    this.xhr.abort();
                }
            }

            // Initializes the XMLHttpRequest object using the URL passed to the constructor.
            _initRequest() {
                const xhr = this.xhr = new XMLHttpRequest();

                xhr.open('POST', '/file/up', true);
                xhr.responseType = 'json';
            }

            // Initializes XMLHttpRequest listeners.
            _initListeners(resolve, reject, file) {
                const xhr = this.xhr;
                const loader = this.loader;
                const genericErrorText = `Couldn't upload file: ${file.name}.`;

                xhr.addEventListener('error', () => reject(genericErrorText));
                xhr.addEventListener('abort', () => reject());
                xhr.addEventListener('load', () => {
                    const response = xhr.response;

                    // This example assumes the XHR server's "response" object will come with
                    // an "error" which has its own "message" that can be passed to reject()
                    // in the upload promise.
                    //
                    // Your integration may handle upload errors in a different way so make sure
                    // it is done properly. The reject() function must be called when the upload fails.
                    if (!response || response.error) {
                        return reject(response && response.error ? response.error.message : genericErrorText);
                    }

                    // If the upload is successful, resolve the upload promise with an object containing
                    // at least the "default" URL, pointing to the image on the server.
                    // This URL will be used to display the image in the content. Learn more in the
                    // UploadAdapter#upload documentation.

                    resolve({
                        default: response.Url
                    });
                });

                // Upload progress when it is supported. The file loader has the #uploadTotal and #uploaded
                // properties which are used e.g. to display the upload progress bar in the editor
                // user interface.
                if (xhr.upload) {
                    xhr.upload.addEventListener('progress', evt => {
                        if (evt.lengthComputable) {
                            loader.uploadTotal = evt.total;
                            loader.uploaded = evt.loaded;
                        }
                    });
                }
            }

            // Prepares the data and sends the request.
            _sendRequest(file) {
                // Prepare the form data.
                const data = new FormData();


                data.append('file', file);

                // Important note: This is the right place to implement security mechanisms
                // like authentication and CSRF protection. For instance, you can use
                // XMLHttpRequest.setRequestHeader() to set the request headers containing
                // the CSRF token generated earlier by your application.

                // Send the request.
                this.xhr.send(data);
            }
        }

        function MyCustomUploadAdapterPlugin(editor) {
            editor.plugins.get('FileRepository').createUploadAdapter = (loader) => {
                // Configure the URL to the upload script in your back-end here!
                return new MyUploadAdapter(loader);
            };
        }

        await loadArticle(hasArticleID);

        /*DecoupledEditor
            .create(document.querySelector('#editor-container'), {
                // toolbar: [ 'heading', '|', 'bold', 'italic', 'link' ]
                extraPlugins: [MyCustomUploadAdapterPlugin],
                language: "zh-cn"
            })
            .then(editor => {

                const toolbarContainer = document.querySelector('#editor-toolbar-container');

                toolbarContainer.prepend(editor.ui.view.toolbar.element);

                window.editor = editor;

                loadArticle(hasArticleID);
            })
            .catch(err => {
                console.error(err.stack);
            });*/


    });


});

main.controller("content_add_gallery_controller", function ($http, $scope, $routeParams, $rootScope, $timeout, $location, Upload) {
    $scope.ContentSubTypes = {};


    $scope.ContentItemID = parseInt($routeParams.ContentItemID);
    //$scope.Article={ContentItemID:$scope.ContentItemID};

    //$scope.ContentSubTypeID;
    $scope.MContentSubTypeID = 0;
    $scope.MContentSubTypeChildID = 0;


    $scope.articles = [];

    $scope.saveArticle = async function () {

        /*const upload_progress = $("#upload_article_images_progress");
        upload_progress.progress({
            duration : 100,
            total    : 100,
            text:{
                active: '{value} of {total} done'
            }
        });

        upload_progress.progress('update progress',50);*/


        for (let i = 0; i < $scope.articles.length; i++) {

            let article = $scope.articles[i];
            if (article.PictureBlob) {

                let p = await new Promise((resolve, reject) => {

                    Upload.upload({url: '/file/up', data: {file: article.PictureBlob}}).then(function (response) {
                        const url = response.data.Path;
                        //$scope.articles.push({Picture:url})
                        resolve(url)

                    }, function (response) {

                    }, function (evt) {

                        //const progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                        //upload_progress.progress('update progress',progress);

                    });

                })
                article.Picture = p
                delete article["PictureBlob"]

            }

            article.ContentSubTypeID = parseInt(article.ContentSubTypeID)

            if (parseInt(article.ContentSubTypeChildID) > 0) {
                article.ContentSubTypeID = parseInt(article.ContentSubTypeChildID)
            }

            let p = await new Promise((resolve, reject) => {
                $http.post("content/save", JSON.stringify(article), {
                    transformRequest: angular.identity,
                    headers: {"Content-Type": "application/json"}
                }).then(function (data, status, headers, config) {
                    resolve()
                });
            })

            $("#upload_article_images_progress").progress('update progress', parseInt((i / $scope.articles.length) * 100))


        }

        $("#upload_article_images_progress").progress('update progress', 100)

        $timeout(function () {
            $scope.articles = [];
        })
    }

    $scope.listContentSubTypes = function () {
        //content/list
        $http.get("content/sub_type/list/all/" + $routeParams.ContentItemID).then(function (value) {

            $scope.ContentSubTypes = value.data.Data;

        });
    }
    $scope.listContentSubTypes();

    $scope.changeContentSubTypes = function (ContentSubTypeID) {
        //$scope.ContentSubTypeChildID=undefined;
        //$scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        //$scope.ContentSubTypeChilds=[];
        //console.log($scope.ContentSubTypeID);

        for (let i = 0; i < $scope.articles.length; i++) {
            $scope.articles[i].ContentSubTypeID = ContentSubTypeID;
            $scope.articles[i].ContentSubTypeChildID = 0;
        }


    }
    $scope.changeContentSubTypeChilds = function () {
        //alert($scope.MContentSubTypeChildID);
        // if($scope.MContentSubTypeChildID){
        //     $scope.Article.ContentSubTypeID=$scope.MContentSubTypeChildID;
        // }else{
        //     $scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        // }

        //$scope.Article.;

        //alert($scope.MContentSubTypeChildID)

        for (let i = 0; i < $scope.articles.length; i++) {
            $scope.articles[i].ContentSubTypeChildID = $scope.MContentSubTypeChildID;
        }

    }


    $scope.changeArticleContentSubTypes = function (m) {


        m.ContentSubTypeChildID = 0;

    }


    $scope.UploadImages = function (progressID, files, errFiles) {

        /*  const upload_progress = $(progressID);
          upload_progress.progress({
              duration : 100,
              total    : 100,
              text:{
                  active: '{value} of {total} done'
              }
          });*/

        //upload_progress.progress('reset');
        //upload_progress.progress('update progress',50);

        if (files && files.length) {
            for (let i = 0; i < files.length; i++) {
//PictureBlob
                console.log(files[i])
                $scope.articles.push({
                    Content: files[i].type,
                    Title: files[i].name,
                    PictureBlob: files[i],
                    ContentItemID: $scope.ContentItemID,
                    ContentSubTypeChildID: $scope.MContentSubTypeChildID,
                    ContentSubTypeID: $scope.MContentSubTypeID
                })
                /*Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    const url = response.data.Data;
                    $scope.articles.push({Picture:url})

                },function (response) {

                },function (evt) {

                    const progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                    upload_progress.progress('update progress',progress);

                });*/

            }
        } else {
            UpImageError(errFiles);
        }
    }


});

main.controller('content_list_controller', function ($http, $scope, $timeout, $routeParams, $document, $interval) {

    $scope.MenuTypes = [];
    $scope.EditMenus = null;
    $scope.CreateMenus = null;
    $scope.selectClassify = null;
    $scope.classifyChild = {};

    $scope.ContentContentSubType = {}
    $scope.ContentSubTypes = {}

    $scope.templateNameObj = {
        "contents": [
            {Key: "news", Label: "新闻中心", SubMenu: true, Content: true}
        ],
        "content": [
            {Key: "services", Label: "服务", SubMenu: true, Content: true},
            {Key: "about", Label: "关于我们", SubMenu: true, Content: true}
        ],
        "index": [
            {Key: "index", Label: "首页", SubMenu: false, Content: false}
        ],
        "gallery": [
            {Key: "gallery", Label: "媒体", SubMenu: true, Content: true}
        ],
        "products": [
            {Key: "products", Label: "产品", SubMenu: false, Content: false}
        ],
        "blog": [
            {Key: "blog", Label: "博客", SubMenu: true, Content: true}
        ],
    };

    $scope.getTemplateNameObj = function (type, templateName) {
        let tns = $scope.templateNameObj[type];
        for (let i = 0; i < tns.length; i++) {
            if (tns[i].Key === templateName) {
                return tns[i]
            }
        }

    }
    $scope.templateNameObjFunc = function (contentTypeID) {

        for (let i = 0; i < $scope.MenuTypes.length; i++) {
            if ($scope.MenuTypes[i].ID === contentTypeID) {
                return $scope.templateNameObj[$scope.MenuTypes[i].Type];

            }
        }

        return [];

    }

    let ActionTarget = {method: 'POST', url: 'menus', title: '添加菜单'};

    $http.get("content/type/list").then(function (value) {

        $scope.MenuTypes = value.data.Data;

    });
    $scope.listMenus = function () {
        //content/list
        $http.get("content/item/list").then(function (value) {

            $scope.MenusList = value.data.Data;

        });
    }
    $scope.listClassify = function () {
        //content/list
        $http.get("content/sub_type/list/" + $scope.EditMenus.ID).then(function (value) {

            $scope.ClassifyList = value.data.Data;

        });
    }
    $scope.listChildClassify = function (ContentItemID, ParentID) {
        //content/list
        $http.get("content/sub_type/child/list/" + ContentItemID + "/" + ParentID).then(function (value) {

            $scope.ClassifyChildList = value.data.Data;

        });
    }
    $scope.saveMenu = async function (menu) {
        return await new Promise((resolve, reject) => {
            $http({
                method: ActionTarget.method,
                url: ActionTarget.url,
                data: JSON.stringify(menu),
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json;charset=utf-8"}
            }).then(function (data) {
                $scope.listMenus();
                alert(data.data.Message);
                resolve(true);

            }).catch((error) => {
                resolve(false);
            });
        });
    }
    $scope.upIndex = function (index) {
        if (index == 0) {
            return
        }
        const current = angular.copy($scope.MenusList[index]);//1
        const targetIndex = (index - 1) <= 0 ? 0 : (index - 1);
        const target = angular.copy($scope.MenusList[targetIndex]);//0

        $scope.changeMenuSort(current, index, targetIndex, target);

    }


    $scope.changeHide = function (m) {

        ActionTarget = {method: 'PUT', url: 'content/item/hide/' + m.ID, title: '修改显示'};
        $http({
            method: ActionTarget.method,
            url: ActionTarget.url,
            data: JSON.stringify({Hide: m.Hide}),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {

            $scope.listMenus();
            //$scope.Menus = null;

        });

    }
    $scope.downIndex = function (index) {

        if ($scope.MenusList.length - 1 === index) {
            return
        }
        const current = angular.copy($scope.MenusList[index]);//1
        const targetIndex = (index + 1) >= $scope.MenusList.length - 1 ? $scope.MenusList.length - 1 : (index + 1);
        const target = angular.copy($scope.MenusList[targetIndex]);//0


        $scope.changeMenuSort(current, index, targetIndex, target);

    }
    $scope.changeMenuSort = function (current, orgIndex, targetIndex, target) {
        $scope.MenusList[targetIndex] = current;
        $scope.MenusList[orgIndex] = target;

        current.Sort = targetIndex;
        target.Sort = orgIndex;

        ActionTarget = {method: 'PUT', url: 'content/item/index/' + target.ID, title: '修改菜单'};
        $http({
            method: ActionTarget.method,
            url: ActionTarget.url,
            data: JSON.stringify(target),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {
            ActionTarget = {method: 'PUT', url: 'content/item/index/' + current.ID, title: '修改菜单'};
            $http({
                method: ActionTarget.method,
                url: ActionTarget.url,
                data: JSON.stringify(current),
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json;charset=utf-8"}
            }).then(function (data) {

                $scope.listMenus();
                //$scope.Menus = null;

            });
        });
    }
    $scope.saveCreateMenu = async function () {
        ActionTarget = {method: 'POST', url: 'content/item/add', title: '添加菜单'};
        await $scope.saveMenu($scope.CreateMenus);
        $scope.CreateMenus = null;
    }
    $scope.saveEditMenu = async function () {
        ActionTarget = {method: 'PUT', url: 'content/item/' + $scope.EditMenus.ID, title: '修改菜单'};
        await $scope.saveMenu($scope.EditMenus);
        $scope.EditMenus = null;
    }
    //{method:'PUT',url:''}


    $scope.changeContentSubTypes = function () {
        $scope.loadContent()
    }
    $scope.changeContentSubTypeChilds = function () {
        $scope.loadContent()
    }
    $scope.loadContent = function () {
        let ContentSubTypeID = 0

        if (isNaN(parseInt($scope.ContentContentSubType.ContentSubTypeID)) === false) {
            ContentSubTypeID = parseInt($scope.ContentContentSubType.ContentSubTypeID)
        }

        if (isNaN(parseInt($scope.ContentContentSubType.ContentSubTypeChildID)) === false) {
            ContentSubTypeID = parseInt($scope.ContentContentSubType.ContentSubTypeChildID)
        }

        if (ContentSubTypeID === 0) {
            //alert("请选择类别")
            //return
        }

        $http.get("content/single/get/" + $scope.ContentContentSubType.ContentItemID + "/" + ContentSubTypeID).then(function (responea) {

            $scope.Article = responea.data.Data || {};

        });
    }

    $scope.deleteArticle = function (article) {
        if (confirm("确定删除？")) {

            const form = {};
            form.ID = article.ID;
            $http.post("content/delete", $.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(async function (data, status, headers, config) {

                alert(data.data.Message);
                $scope.loadContent();

            });

        }
    }

    $scope.gotoContent = function () {
        $timeout(function () {
            let redirect = "#!/content?Single=true&ContentItemID=" + $scope.ContentContentSubType.ContentItemID + "&ContentSubTypeID=" + $scope.ContentContentSubType.ContentSubTypeID + "&ContentSubTypeChildID=" + $scope.ContentContentSubType.ContentSubTypeChildID
            //$scope.ContentSubTypes = {};
            //$scope.ContentContentSubType = {};
            window.open(redirect)
            //$("#contentContentSubTypeDialog").modal("hide");
        })

    }

    async function loadAllContentSubType(ContentItemID) {
        await new Promise((resolve, reject) => {
            $http.get("content/sub_type/list/all/" + ContentItemID).then(function (value) {

                $scope.ContentSubTypes = value.data.Data || [];

                resolve();
            });
        });
    }

    $scope.showContentContentSubTypeDialog = async function (ContentItemID) {
        $scope.ContentContentSubType.ContentItemID = ContentItemID


        await loadAllContentSubType(ContentItemID);


        $scope.loadContent()


        $("#contentContentSubTypeDialog").modal({
            centered: false, closable: false, allowMultiple: false,
            onDeny: function (e) {
                $timeout(function () {
                    $scope.ContentSubTypes = {};
                    $scope.ContentContentSubType = {};
                })
                return true
            },
            onApprove: function (e) {

                if ($scope.ContentContentSubType.ContentSubTypeID == undefined) {
                    alert("请选择类别")
                    return false
                }
                if ($scope.ContentContentSubType.ContentSubTypeChildID == undefined) {
                    $scope.ContentContentSubType.ContentSubTypeChildID = 0
                }

                if (isNaN(parseInt($scope.ContentContentSubType.ContentSubTypeID)) || isNaN(parseInt($scope.ContentContentSubType.ContentSubTypeChildID))) {
                    alert("类别选择错误")
                    return
                }

                $timeout(function () {
                    //let redirect = "#!/" + m.Type + "?ContentItemID=" + m.ID + "&ContentSubTypeID=" + $scope.ContentContentSubType.ContentSubTypeID + "&ContentSubTypeChildID=" + $scope.ContentContentSubType.ContentSubTypeChildID
                    $scope.ContentSubTypes = {};
                    $scope.ContentContentSubType = {};
                    //window.location.href = redirect
                })

                return true
            }
        }).modal("show");

    }
    $scope.editMenus = function (m) {
        //$scope.selectClassify = null;
        //$scope.classifyChild = null;
        $scope.EditMenus = m;
        $scope.EditMenus.Template = $scope.getTemplateNameObj($scope.EditMenus.Type, $scope.EditMenus.TemplateName);

        $scope.classify = {ContentItemID: $scope.EditMenus.ID};
        $("#editMenus").modal({
            centered: false, closable: false, allowMultiple: false,
            onApprove: function (e) {

                $timeout(function () {
                    $scope.selectClassify = null;
                    $scope.classifyChild = {};
                });

                return true
            }
        }).modal("show");

        $scope.listClassify();
    }
    $scope.deleteMenus = function (ID) {
        if (confirm("确定要删除？") == false) {
            return
        }
        $http.delete("content/item/" + ID, {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {

            alert(data.data.Message);

            $scope.listMenus();

            // $scope.Menus = null;

        });
    }
    $scope.listMenus();


    $scope.classify = null;

    $scope.ActionClassifyTarget = {method: 'POST', url: 'content/sub_type', title: '添加分类'};

    $scope.deleteClassify = function (m) {

        if (confirm("确定要删除？") == false) {
            return
        }

        $http.delete("content/sub_type/" + m.ID, {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {

            alert(data.data.Message);
            $scope.listClassify();

        });

    }
    $scope.editClassify = function (m) {

        $scope.classify = m;
        $scope.ActionClassifyTarget = {method: 'PUT', url: 'content/sub_type/' + m.ID, title: '修改分类'};
        //$scope.saveClassify();
    }
    $scope.saveClassify = function () {

        $http({
            method: $scope.ActionClassifyTarget.method,
            url: $scope.ActionClassifyTarget.url,
            data: JSON.stringify($scope.classify),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {

            $scope.listClassify();
            $scope.classify.Name = '';
            $scope.classify.ID = null;
            $scope.ActionClassifyTarget = {method: 'POST', url: 'content/sub_type', title: '添加分类'};
            alert(data.data.Message);
        });


    }


    $scope.selectClassifyFunc = function (m) {
        //$scope.ActionClassifyTarget={method:'PUT',url:'content_sub_type/'+m.ID,title:'修改分类'};
        $scope.selectClassify = m;
        $scope.listChildClassify($scope.selectClassify.ContentItemID, $scope.selectClassify.ID);
    }


    $scope.ActionClassifyChildTarget = {method: 'POST', url: 'content/sub_type', title: '添加子分类'};

    //saveClassifyChild
    $scope.saveClassifyChild = function () {

        if (!$scope.selectClassify) {
            alert("请选择父类");
            return
        }
        if (!$scope.EditMenus) {
            alert("请菜单");
            return
        }


        $scope.classifyChild.ParentContentSubTypeID = $scope.selectClassify.ID;
        $scope.classifyChild.ContentItemID = $scope.EditMenus.ID;


        $http({
            method: $scope.ActionClassifyChildTarget.method,
            url: $scope.ActionClassifyChildTarget.url,
            data: JSON.stringify($scope.classifyChild),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {
            alert(data.data.Message);
            $scope.listChildClassify($scope.selectClassify.ContentItemID, $scope.selectClassify.ID);

            //$scope.classifyChild.Name = '';
            //$scope.classifyChild.ID = null;
            $scope.classifyChild = {};
            $scope.ActionClassifyChildTarget = {method: 'POST', url: 'content/sub_type', title: '添加分类'};
        });


    }
    $scope.deleteClassifyChild = function (m) {
        if (confirm("确定要删除？") == false) {
            return
        }
        $http.delete("content/sub_type/" + m.ID, {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json;charset=utf-8"}
        }).then(function (data) {

            alert(data.data.Message);
            $scope.listChildClassify($scope.selectClassify.ContentItemID, $scope.selectClassify.ID);

        });

    }
    $scope.editClassifyChild = function (m) {
        $scope.classifyChild = m;
        $scope.ActionClassifyChildTarget = {method: 'PUT', url: 'content/sub_type/' + m.ID, title: '修改分类'};
        //$scope.saveClassify();
    }


});
