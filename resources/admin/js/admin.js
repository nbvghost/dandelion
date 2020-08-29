var main = angular.module("managerApp", ['ngRoute',"ngMessages","ngFileUpload"]);
//main.config(function ($interpolateProvider){$interpolateProvider.startSymbol("@{").endSymbol("}@");});

main.config(function($routeProvider, $locationProvider,$provide,$httpProvider,$httpParamSerializerJQLikeProvider,$interpolateProvider) {
    $interpolateProvider.startSymbol("@{").endSymbol("}@");
    $httpProvider.defaults.transformRequest.unshift($httpParamSerializerJQLikeProvider.$get());
    $httpProvider.defaults.headers.post={'Content-Type':'application/x-www-form-urlencoded;charset=UTF-8'};

    /*$provide.factory('httpInterceptor', function($q,$rootScope,$timeout) {
        return {
            'request': function(config) {
                $rootScope.progressbar.start();
                $timeout($rootScope.progressbar.complete(), 10000);
                return config;
            },
            'requestError': function(rejection) {
                $rootScope.progressbar.complete();
                Messager(rejection.status+":"+rejection.statusText);
                return rejection;
            },
            'response': function(response) {
                $rootScope.progressbar.complete();
                return response;
            },
            'responseError': function(rejection) {
                $rootScope.progressbar.complete();
                Messager(rejection.status+":"+rejection.statusText);
                return rejection;
            }
        };
    });

    $httpProvider.interceptors.push("httpInterceptor");*/



    $routeProvider.when("/", {
        templateUrl: "main",
        controller: "main_controller"
    });

    $routeProvider.when("/add_goods", {
        templateUrl: "add_goods",
        controller: "add_goods_controller"
    });

    $routeProvider.when("/goods_list", {
        templateUrl: "goods_list",
        controller: "goods_list_controller"
    });

    $routeProvider.when("/store_stock_manager", {
        templateUrl: "store_stock_manager",
        controller: "store_stock_manager_controller"
    });

    $routeProvider.when("/score_goods_list", {
        templateUrl: "score_goods_list",
        controller: "score_goods_list_controller"
    });
    $routeProvider.when("/voucher_list", {
        templateUrl: "voucher_list",
        controller: "voucher_list_controller"
    });
    $routeProvider.when("/fullcut_list", {
        templateUrl: "fullcut_list",
        controller: "fullcut_list_controller"
    });

    $routeProvider.when("/timesell_list", {
        templateUrl: "timesell_list",
        controller: "timesell_list_controller"
    });
    $routeProvider.when("/add_timesell", {
        templateUrl: "add_timesell",
        controller: "add_timesell_controller"
    });
    $routeProvider.when("/timesell_manager", {
        templateUrl: "timesell_manager",
        controller: "timesell_manager_controller"
    });
    $routeProvider.when("/collage_manager", {
        templateUrl: "collage_manager",
        controller: "collage_manager_controller"
    });
    $routeProvider.when("/collage_list", {
        templateUrl: "collage_list",
        controller: "collage_list_controller"
    });
    $routeProvider.when("/add_collage", {
        templateUrl: "add_collage",
        controller: "add_collage_controller"
    });

    $routeProvider.when("/add_store", {
        templateUrl: "add_store",
        controller: "add_store_controller"
    });
    $routeProvider.when("/store_list", {
        templateUrl: "store_list",
        controller: "store_list_controller"
    });
    $routeProvider.when("/admin_list", {
        templateUrl: "admin_list",
        controller: "admin_list_controller"
    });
    $routeProvider.when("/express", {
        templateUrl: "express",
        controller: "express_controller"
    });
    $routeProvider.when("/add_express", {
        templateUrl: "add_express",
        controller: "add_express_controller"
    });
    $routeProvider.when("/order_list", {
        templateUrl: "order_list",
        controller: "order_list_controller"
    });

    $routeProvider.when("/user_setup", {
        templateUrl: "user_setup",
        controller: "user_setup_controller"
    });

    $routeProvider.when("/view_situation", {
        templateUrl: "view_situation",
        controller: "view_situation_controller"
    });
    $routeProvider.when("/carditem_list", {
        templateUrl: "carditem_list",
        controller: "carditem_list_controller"
    });
    $routeProvider.when("/store_situation", {
        templateUrl: "store_situation",
        controller: "store_situation_controller"
    });

    $routeProvider.when("/content_list", {
        templateUrl: "content_templets/content_list_templet",
        controller: "content_list_controller"
    });
    $routeProvider.when("/add_articles", {
        templateUrl: "content_templets/add_articles_templet",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/add_gallery", {
        templateUrl: "content_templets/add_gallery_templet",
        controller: "content_add_gallery_controller"
    });
    $routeProvider.when("/articles", {
        templateUrl: "content_templets/articles_templet",
        controller: "content_articles_controller"
    });
    $routeProvider.when("/gallery", {
        templateUrl: "content_templets/articles_templet",
        controller: "content_articles_controller"
    });
    $routeProvider.when("/article", {
        templateUrl: "content_templets/add_article_templet",
        controller: "content_add_article_controller"
    });






    //$locationProvider.html5Mode(true)//.hashPrefix('#!');
    
});


main.controller("content_articles_controller", function ($http, $scope, $routeParams, $rootScope,$timeout,$location,Upload) {
    $scope.ContentSubTypes=[];
    $scope.ContentSubTypeChilds=[];

    $scope.ContentSubTypeID;
    $scope.MContentSubTypeID;
    $scope.MContentSubTypeChildID;

    $scope.ContentID=$routeParams.ContentID;
    $scope.Type=$routeParams.Type;

    var table;

    $scope.listContentSubTypes = function(){
        //content/list
        $http.get("content/sub_type/list/"+$routeParams.ContentID).then(function (value){

            $scope.ContentSubTypes = value.data.Data;

        });
    }
    $scope.listContentSubTypes();


    $scope.listContentSubTypeChilds = function(ContentID,ParentContentSubTypeID){
        //content/list
        $http.get("content/sub_type/child/list/"+ContentID+"/"+ParentContentSubTypeID).then(function (value){

            $scope.ContentSubTypeChilds = value.data.Data;

        });
    }


    $scope.changeContentSubTypeChilds = function(){
        $scope.ContentSubTypeID=$scope.MContentSubTypeChildID;
        table.ajax.reload();

    }
    $scope.changeContentSubTypes = function(){
        $scope.ContentSubTypeID=$scope.MContentSubTypeID;
        $scope.ContentSubTypeChilds=[];
        if($scope.ContentSubTypeID){
            $scope.listContentSubTypeChilds($scope.ContentID,$scope.ContentSubTypeID);
        }
        table.ajax.reload();
    }

    $timeout(function () {

        table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Author"},
                {data:"Look"},
                {data:"ContentID",visible:false},
                {data:"ContentSubTypeID",visible:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "order":[[1,"asc"]],
            "ajax": {
                "url": "article/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    d.columns[4].search.value=parseInt($scope.ContentID).toString();


                    if($scope.ContentSubTypeID){
                        d.columns[5].search.value=parseInt($scope.ContentSubTypeID).toString();
                    }
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());
            window.location.href="#!/add_articles?ContentID="+row.data().ContentID+"&ID="+row.data().ID;
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    var form ={};
                    form.ID =row.data().ID;
                    $http.post("article/delete",$.param(form), {
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
main.controller("content_add_articles_controller", function ($http, $scope, $routeParams, $rootScope,$timeout,$location,Upload) {
    $scope.ContentSubTypes=[];
    $scope.ContentSubTypeChilds=[];



    $scope.Article={ContentID:$routeParams.ContentID};

    //$scope.ContentSubTypeID;
    $scope.MContentSubTypeID;
    $scope.MContentSubTypeChildID;


    console.log($location)


    $scope.saveArticle = function(){

        //$scope.ContentSubTypeID;
        //$scope.ContentSubTypeChildID;
        //console.log(quill.container.firstChild.innerHTML)
        $scope.Article.ContentID=parseInt($routeParams.ContentID);

        $scope.Article.Content=quill.container.firstChild.innerHTML;

        if(!$scope.Article.ContentSubTypeID){
            //$scope.Article.ContentSubTypeID=$scope.ContentSubTypeID;
            //$scope.Article.ContentSubTypeChildID=$scope.ContentSubTypeChildID;
        //}else{
            alert("请选择分类");
            return
        }
        $http.post("article/save",JSON.stringify($scope.Article), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {
            console.log(data);
            alert(data.data.Message);
            if(data.data.Success){
                window.history.back();
            }
        });
    }

    $scope.listContentSubTypes = function(){
        //content/list
        $http.get("content/sub_type/list/"+$routeParams.ContentID).then(function (value){

            $scope.ContentSubTypes = value.data.Data;

        });
    }
    $scope.listContentSubTypes();


    $scope.changeContentSubTypes = function(){
        //$scope.ContentSubTypeChildID=undefined;
        $scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        $scope.ContentSubTypeChilds=[];
        console.log($scope.ContentSubTypeID);
        $scope.listContentSubTypeChilds($routeParams.ContentID,$scope.Article.ContentSubTypeID);
    }
    $scope.changeContentSubTypeChilds = function(){
        //alert($scope.MContentSubTypeChildID);
        if($scope.MContentSubTypeChildID){
            $scope.Article.ContentSubTypeID=$scope.MContentSubTypeChildID;
        }else{
            $scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        }

        //$scope.Article.;
    }

    //content_sub_type/child/list/:ParentContentSubTypeID

    $scope.listContentSubTypeChilds = function(ContentID,ParentContentSubTypeID){
        //content/list
        $http.get("content/sub_type/child/list/"+ContentID+"/"+ParentContentSubTypeID).then(function (value){

            $scope.ContentSubTypeChilds = value.data.Data;

        });
    }


    $scope.UploadPictureImage = function (file, errFiles) {
        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                var url =response.data.Data;
                $scope.Article.Picture=url;
            }, function (response) {
                console.log(response);
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //upload_progress.progress('update progress',progress);
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            if(errFiles.length>0){
                alert(JSON.stringify(errFiles));
            }

        }
    }

    $scope.EditImages=[];
    $scope.UploadImages = function (progressID,files, errFiles) {

        var upload_progress = $(progressID);
        upload_progress.progress({
            duration : 100,
            total    : 100,
            text:{
                active: '{value} of {total} done'
            }
        });

        upload_progress.progress('reset');
        //upload_progress.progress('update progress',50);

        if (files && files.length) {
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Data;

                    if($scope.EditImages.indexOf(url)==-1){
                        $scope.EditImages.push(url);
                    }

                },function (response) {

                },function (evt) {

                    var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                    upload_progress.progress('update progress',progress);

                });
            }
        }else{
            UpImageError(errFiles);
        }
    }
    var quill;
    $timeout(function () {

        var Inline = Quill.import('blots/inline');
        var Block = Quill.import('blots/block');
        var BlockEmbed = Quill.import('blots/block/embed');

        class BoldBlot extends Inline { }
        BoldBlot.blotName = 'bold';
        BoldBlot.tagName = 'strong';

        class ItalicBlot extends Inline { }
        ItalicBlot.blotName = 'italic';
        ItalicBlot.tagName = 'em';

        class LinkBlot extends Inline {
            static create(url) {
                var node = super.create();
                node.setAttribute('href', url);
                node.setAttribute('target', '_blank');
                return node;
            }

            static formats(node) {
                return node.getAttribute('href');
            }
        }
        LinkBlot.blotName = 'link';
        LinkBlot.tagName = 'a';

        class BlockquoteBlot extends Block { }
        BlockquoteBlot.blotName = 'blockquote';
        BlockquoteBlot.tagName = 'blockquote';

        class HeaderBlot extends Block {
            static formats(node) {
                return HeaderBlot.tagName.indexOf(node.tagName) + 1;
            }
        }
        HeaderBlot.blotName = 'header';
        HeaderBlot.tagName = ['H1', 'H2'];

        class DividerBlot extends BlockEmbed { }
        DividerBlot.blotName = 'divider';
        DividerBlot.tagName = 'hr';

        class ImageBlot extends BlockEmbed {
            static create(value) {
                var node = super.create();
                node.setAttribute('alt', value.alt);
                node.setAttribute('src', value.url);
                return node;
            }

            static value(node) {
                return {
                    alt: node.getAttribute('alt'),
                    url: node.getAttribute('src')
                };
            }
        }
        ImageBlot.blotName = 'image';
        ImageBlot.tagName = 'img';

        class VideoBlot extends BlockEmbed {
            static create(url) {
                var node = super.create();
                node.setAttribute('src', url);
                node.setAttribute('frameborder', '0');
                node.setAttribute('allowfullscreen', true);
                return node;
            }

            static formats(node) {
                var format = {};
                if (node.hasAttribute('height')) {
                    format.height = node.getAttribute('height');
                }
                if (node.hasAttribute('width')) {
                    format.width = node.getAttribute('width');
                }
                return format;
            }

            static value(node) {
                return node.getAttribute('src');
            }

            format(name, value) {
                if (name === 'height' || name === 'width') {
                    if (value) {
                        this.domNode.setAttribute(name, value);
                    } else {
                        this.domNode.removeAttribute(name, value);
                    }
                } else {
                    super.format(name, value);
                }
            }
        }
        VideoBlot.blotName = 'video';
        VideoBlot.tagName = 'iframe';

        Quill.register(BoldBlot);
        Quill.register(ItalicBlot);
        Quill.register(LinkBlot);
        Quill.register(BlockquoteBlot);
        Quill.register(HeaderBlot);
        Quill.register(DividerBlot);
        Quill.register(ImageBlot);
        Quill.register(VideoBlot);

        quill = new Quill('#editor-container', {
            modules: {
                formula: true,
                syntax: true,
                toolbar: '#toolbar-container'
            },
            placeholder: 'Compose an epic...',
            theme: 'snow'
        });

        if($routeParams.ID){
            //article/get/:ID
            $http.get("article/multi/get/"+$routeParams.ID).then(function (responea){

                $scope.Article = responea.data.Data;
                quill.clipboard.dangerouslyPasteHTML($scope.Article.Content);

                //content_sub_type
                $http.get("content/sub_type/"+$scope.Article.ContentSubTypeID).then(function (responeb){
                    var ContentSubType = responeb.data.Data.ContentSubType;
                    var ParentContentSubType = responeb.data.Data.ParentContentSubType;


                    $timeout(function () {
                        if(ParentContentSubType.ID>0){
                            $scope.MContentSubTypeID=ParentContentSubType.ID;
                            $scope.MContentSubTypeChildID=ContentSubType.ID;
                            $scope.listContentSubTypeChilds($scope.Article.ContentID,$scope.MContentSubTypeID);
                        }else{
                            $scope.MContentSubTypeID=ContentSubType.ID;
                            $scope.MContentSubTypeChildID=ParentContentSubType.ID;
                            $scope.listContentSubTypeChilds($scope.Article.ContentID,$scope.MContentSubTypeID);
                        }

                    });



                });






            });
        }

        quill.getModule("toolbar").addHandler("image", function (e) {

            //var baseUrl ="//"+$location.host()+":"+$location.port();

            $("#SelectImageModal").modal({onApprove:function (e) {


                    if($scope.EditImages.length>0){


                        for(var ii=0;ii<$scope.EditImages.length;ii++){

                            var range = quill.getSelection(true);
                            quill.insertText(range.index, '\n', Quill.sources.USER);
                            quill.insertEmbed(range.index + 1, 'image', {
                                alt: '软件定制开发，QQ/微信：274455411',
                                url: $scope.EditImages[ii]
                            }, Quill.sources.USER);
                            quill.setSelection(range.index + 2, Quill.sources.SILENT);
                        }


                            /*for(var ii=0;ii<$scope.EditImages.length;ii++){

                                quill.insertEmbed(range.index, 'image',);
                                range = quill.getSelection();
                            }*/

                            $scope.$apply(function () {
                                $scope.EditImages=[];
                            });
                            return true;



                    }else{
                        return false;
                    }
                },closable:false}).modal("show");
        });

    });


});
main.controller("content_add_gallery_controller", function ($http, $scope, $routeParams, $rootScope,$timeout,$location,Upload) {
    $scope.ContentSubTypes=[];
    $scope.ContentSubTypeChilds={};


    $scope.ContentID=parseInt($routeParams.ContentID);
    //$scope.Article={ContentID:$scope.ContentID};

    //$scope.ContentSubTypeID;
    $scope.MContentSubTypeID=0;
    $scope.MContentSubTypeChildID=0;


    $scope.articles=[];

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
        $("#upload_article_images_progress").progress('update progress', 50)

        for (let i = 0; i < $scope.articles.length; i++) {

            let article = $scope.articles[i];
            if (article.PictureBlob) {

                let p = await new Promise((resolve, reject) => {

                    Upload.upload({url: '/file/up', data: {file: article.PictureBlob}}).then(function (response) {
                        const url = response.data.Data;
                        //$scope.articles.push({Picture:url})
                        resolve(url)

                    }, function (response) {

                    }, function (evt) {

                        //const progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                        //upload_progress.progress('update progress',progress);

                    });

                })
                console.log(p)
                article.Picture =p
                delete article["PictureBlob"]

            }


            $http.post("article/save",JSON.stringify(article), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json"}
            }).then(function (data, status, headers, config) {
                console.log(data);
            });


        }


    }

    $scope.listContentSubTypes = function(){
        //content/list
        $http.get("content/sub_type/list/"+$routeParams.ContentID).then(function (value){

            $scope.ContentSubTypes = value.data.Data;

        });
    }
    $scope.listContentSubTypes();


    $scope.changeArticleContentSubTypes = function(m){


        m.ContentSubTypeChildID=0;

        $scope.listContentSubTypeChilds($scope.ContentID,m.ContentSubTypeID);
    }

    $scope.changeContentSubTypes = function(ContentSubTypeID){
        //$scope.ContentSubTypeChildID=undefined;
        //$scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        //$scope.ContentSubTypeChilds=[];
        //console.log($scope.ContentSubTypeID);

        for(let i=0;i<$scope.articles.length;i++){
            $scope.articles[i].ContentSubTypeID=ContentSubTypeID;
            $scope.articles[i].ContentSubTypeChildID=0;
        }

        $scope.listContentSubTypeChilds($routeParams.ContentID,ContentSubTypeID);
    }
    $scope.changeContentSubTypeChilds = function(){
        //alert($scope.MContentSubTypeChildID);
        // if($scope.MContentSubTypeChildID){
        //     $scope.Article.ContentSubTypeID=$scope.MContentSubTypeChildID;
        // }else{
        //     $scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        // }

        //$scope.Article.;

        //alert($scope.MContentSubTypeChildID)

        for(let i=0;i<$scope.articles.length;i++){
            $scope.articles[i].ContentSubTypeChildID=$scope.MContentSubTypeChildID;
        }

    }

    //content_sub_type/child/list/:ParentContentSubTypeID

    $scope.listContentSubTypeChilds = function(ContentID,ParentContentSubTypeID){
        //content/list
        $http.get("content/sub_type/child/list/"+ContentID+"/"+ParentContentSubTypeID).then(function (value){

            $scope.ContentSubTypeChilds[ParentContentSubTypeID] = value.data.Data;

        });
    }




    $scope.UploadImages = function (progressID,files, errFiles) {

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
                $scope.articles.push({PictureBlob:files[i],ContentID:$scope.ContentID,ContentSubTypeChildID:$scope.MContentSubTypeChildID,ContentSubTypeID:$scope.MContentSubTypeID})
                /*Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    const url = response.data.Data;
                    $scope.articles.push({Picture:url})

                },function (response) {

                },function (evt) {

                    const progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                    upload_progress.progress('update progress',progress);

                });*/

            }
        }else{
            UpImageError(errFiles);
        }
    }


});
main.controller("content_add_article_controller", function ($http, $scope, $routeParams, $rootScope,$timeout,$location,Upload) {


    $scope.ContentID = $routeParams.ContentID
    //$scope.ID = $routeParams.ID

    $scope.Article={ContentID:$scope.ContentID};





    $scope.saveArticle = function(){

        //$scope.ContentSubTypeID;
        //$scope.ContentSubTypeChildID;
        //console.log(quill.container.firstChild.innerHTML)
        $scope.Article.ContentID=parseInt($scope.ContentID);

        $scope.Article.Content=quill.container.firstChild.innerHTML;
        $scope.Article.ContentSubTypeID=0

        $http.post("article/save",JSON.stringify($scope.Article), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {
            console.log(data);
            alert(data.data.Message);
            if(data.data.Success){
                window.history.back();
            }
        });
    }


    $scope.UploadPictureImage = function (file, errFiles) {
        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                var url =response.data.Data;
                $scope.Article.Picture=url;
            }, function (response) {
                console.log(response);
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //upload_progress.progress('update progress',progress);
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            if(errFiles.length>0){
                alert(JSON.stringify(errFiles));
            }

        }
    }

    $scope.EditImages=[];
    $scope.UploadImages = function (progressID,files, errFiles) {

        var upload_progress = $(progressID);
        upload_progress.progress({
            duration : 100,
            total    : 100,
            text:{
                active: '{value} of {total} done'
            }
        });

        upload_progress.progress('reset');
        //upload_progress.progress('update progress',50);

        if (files && files.length) {
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Data;

                    if($scope.EditImages.indexOf(url)==-1){
                        $scope.EditImages.push(url);
                    }

                },function (response) {

                },function (evt) {

                    var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                    upload_progress.progress('update progress',progress);

                });
            }
        }else{
            UpImageError(errFiles);
        }
    }
    var quill;
    $timeout(function () {

        var Inline = Quill.import('blots/inline');
        var Block = Quill.import('blots/block');
        var BlockEmbed = Quill.import('blots/block/embed');

        class BoldBlot extends Inline { }
        BoldBlot.blotName = 'bold';
        BoldBlot.tagName = 'strong';

        class ItalicBlot extends Inline { }
        ItalicBlot.blotName = 'italic';
        ItalicBlot.tagName = 'em';

        class LinkBlot extends Inline {
            static create(url) {
                var node = super.create();
                node.setAttribute('href', url);
                node.setAttribute('target', '_blank');
                return node;
            }

            static formats(node) {
                return node.getAttribute('href');
            }
        }
        LinkBlot.blotName = 'link';
        LinkBlot.tagName = 'a';

        class BlockquoteBlot extends Block { }
        BlockquoteBlot.blotName = 'blockquote';
        BlockquoteBlot.tagName = 'blockquote';

        class HeaderBlot extends Block {
            static formats(node) {
                return HeaderBlot.tagName.indexOf(node.tagName) + 1;
            }
        }
        HeaderBlot.blotName = 'header';
        HeaderBlot.tagName = ['H1', 'H2'];

        class DividerBlot extends BlockEmbed { }
        DividerBlot.blotName = 'divider';
        DividerBlot.tagName = 'hr';

        class ImageBlot extends BlockEmbed {
            static create(value) {
                var node = super.create();
                node.setAttribute('alt', value.alt);
                node.setAttribute('src', value.url);
                return node;
            }

            static value(node) {
                return {
                    alt: node.getAttribute('alt'),
                    url: node.getAttribute('src')
                };
            }
        }
        ImageBlot.blotName = 'image';
        ImageBlot.tagName = 'img';

        class VideoBlot extends BlockEmbed {
            static create(url) {
                var node = super.create();
                node.setAttribute('src', url);
                node.setAttribute('frameborder', '0');
                node.setAttribute('allowfullscreen', true);
                return node;
            }

            static formats(node) {
                var format = {};
                if (node.hasAttribute('height')) {
                    format.height = node.getAttribute('height');
                }
                if (node.hasAttribute('width')) {
                    format.width = node.getAttribute('width');
                }
                return format;
            }

            static value(node) {
                return node.getAttribute('src');
            }

            format(name, value) {
                if (name === 'height' || name === 'width') {
                    if (value) {
                        this.domNode.setAttribute(name, value);
                    } else {
                        this.domNode.removeAttribute(name, value);
                    }
                } else {
                    super.format(name, value);
                }
            }
        }
        VideoBlot.blotName = 'video';
        VideoBlot.tagName = 'iframe';

        Quill.register(BoldBlot);
        Quill.register(ItalicBlot);
        Quill.register(LinkBlot);
        Quill.register(BlockquoteBlot);
        Quill.register(HeaderBlot);
        Quill.register(DividerBlot);
        Quill.register(ImageBlot);
        Quill.register(VideoBlot);

        quill = new Quill('#editor-container', {
            modules: {
                formula: true,
                syntax: true,
                toolbar: '#toolbar-container'
            },
            placeholder: 'Compose an epic...',
            theme: 'snow'
        });

        if($scope.ContentID){
            //article/get/:ID
            $http.get("article/single/get/"+$scope.ContentID).then(function (responea){

                $scope.Article = responea.data.Data;
                quill.clipboard.dangerouslyPasteHTML($scope.Article.Content);

            });
        }

        quill.getModule("toolbar").addHandler("image", function (e) {

            //var baseUrl ="//"+$location.host()+":"+$location.port();

            $("#SelectImageModal").modal({onApprove:function (e) {


                    if($scope.EditImages.length>0){


                        for(var ii=0;ii<$scope.EditImages.length;ii++){

                            var range = quill.getSelection(true);
                            quill.insertText(range.index, '\n', Quill.sources.USER);
                            quill.insertEmbed(range.index + 1, 'image', {
                                alt: '软件定制开发，QQ/微信：274455411',
                                url: $scope.EditImages[ii]
                            }, Quill.sources.USER);
                            quill.setSelection(range.index + 2, Quill.sources.SILENT);
                        }


                            /*for(var ii=0;ii<$scope.EditImages.length;ii++){

                                quill.insertEmbed(range.index, 'image',);
                                range = quill.getSelection();
                            }*/

                            $scope.$apply(function () {
                                $scope.EditImages=[];
                            });
                            return true;



                    }else{
                        return false;
                    }
                },closable:false}).modal("show");
        });

    });


});
main.controller('content_list_controller', function ($http, $scope, $rootScope, $routeParams,$document,$interval) {

    $scope.MenuTypes=[];
    $scope.Menus;

    $scope.templateNameObj = {
        "articles":[
            {Key:"services",Label:"服务",SubMenu:true,Content:true}
            ],
        "article":[
            {Key:"about",Label:"关于我们",SubMenu:false,Content:true}
        ],
        "index":[
            {Key:"index",Label:"首页",SubMenu:false,Content:false}
        ],
        "gallery":[
            {Key:"gallery",Label:"媒体",SubMenu:true,Content:true}
        ],
        "products":[
            {Key:"products",Label:"产品",SubMenu:false,Content:false}
        ],
    };

    $scope.getTemplateNameObj = function(type,templateName){
        let tns =$scope.templateNameObj[type];
        for(let i=0;i<tns.length;i++){
            if(tns[i].Key===templateName){
                return tns[i]
            }
        }

    }
    $scope.templateNameObjFunc = function(contentTypeID){

        for(let i=0;i<$scope.MenuTypes.length;i++){
            if($scope.MenuTypes[i].ID===contentTypeID){
                return $scope.templateNameObj[$scope.MenuTypes[i].Type];

            }
        }

        return [];

    }
    let ActionTarget = {method: 'POST', url: 'menus', title: '添加菜单'};

    $http.get("content/type/list").then(function (value){

        $scope.MenuTypes = value.data.Data;

    });
    $scope.listMenus = function(){
        //content/list
        $http.get("content/list").then(function (value){

            $scope.MenusList = value.data.Data;

        });
    }
    $scope.listClassify = function(){
        //content/list
        $http.get("content/sub_type/list/"+$scope.Menus.ID).then(function (value){

            $scope.ClassifyList = value.data.Data;

        });
    }
    $scope.listChildClassify = function(ContentID,ParentID){
        //content/list
        $http.get("content/sub_type/child/list/"+ContentID+"/"+ParentID).then(function (value){

            $scope.ClassifyChildList = value.data.Data;

        });
    }
    $scope.saveMenu = function(){

        $http({
            method:ActionTarget.method,
            url:ActionTarget.url,
            data:JSON.stringify($scope.Menus),
            transformRequest:angular.identity,
            headers:{"Content-Type":"application/json;charset=utf-8"}
        }).then(function(data){


            $scope.listMenus();
            $scope.Menus=null;
            alert(data.data.Message);

        });

    }
    $scope.upIndex = function(index){
        if(index==0){
            return
        }
        var current = angular.copy($scope.MenusList[index]);//1
        var targetIndex = (index-1)<=0?0:(index-1);
        var target = angular.copy($scope.MenusList[targetIndex]);//0


        $scope.MenusList[targetIndex]=current;
        $scope.MenusList[index]=target;

        current.Sort=targetIndex;
        target.Sort=index;

        ActionTarget={method:'PUT',url:'content/index/'+target.ID,title:'修改菜单'};
        $http({
            method:ActionTarget.method,
            url:ActionTarget.url,
            data:JSON.stringify(target),
            transformRequest:angular.identity,
            headers:{"Content-Type":"application/json;charset=utf-8"}
        }).then(function(data){
            ActionTarget={method:'PUT',url:'content/index/'+current.ID,title:'修改菜单'};
            $http({
                method:ActionTarget.method,
                url:ActionTarget.url,
                data:JSON.stringify(current),
                transformRequest:angular.identity,
                headers:{"Content-Type":"application/json;charset=utf-8"}
            }).then(function(data){

                $scope.listMenus();
                $scope.Menus=null;

            });
        });




    }



    $scope.changeHide = function(m){

        ActionTarget={method:'PUT',url:'content/hide/'+m.ID,title:'修改显示'};
        $http({
            method:ActionTarget.method,
            url:ActionTarget.url,
            data:JSON.stringify({Hide:m.Hide}),
            transformRequest:angular.identity,
            headers:{"Content-Type":"application/json;charset=utf-8"}
        }).then(function(data){

            $scope.listMenus();
            $scope.Menus=null;

        });

    }
    $scope.downIndex = function(index){


        if($scope.MenusList.length-1==index){
            return
        }
        var current = angular.copy($scope.MenusList[index]);//1
        var targetIndex = (index+1)>=$scope.MenusList.length-1?$scope.MenusList.length-1:(index+1);
        var target = angular.copy($scope.MenusList[targetIndex]);//0


        $scope.MenusList[targetIndex]=current;
        $scope.MenusList[index]=target;

        current.Sort=targetIndex;
        target.Sort=index;

        ActionTarget={method:'PUT',url:'content/index/'+target.ID,title:'修改菜单'};
        $http({
            method:ActionTarget.method,
            url:ActionTarget.url,
            data:JSON.stringify(target),
            transformRequest:angular.identity,
            headers:{"Content-Type":"application/json;charset=utf-8"}
        }).then(function(data){
            ActionTarget={method:'PUT',url:'content/index/'+current.ID,title:'修改菜单'};
            $http({
                method:ActionTarget.method,
                url:ActionTarget.url,
                data:JSON.stringify(current),
                transformRequest:angular.identity,
                headers:{"Content-Type":"application/json;charset=utf-8"}
            }).then(function(data){

                $scope.listMenus();
                $scope.Menus=null;

            });
        });

    }
    $scope.saveMenuInline = function(){
        ActionTarget={method:'POST',url:'content/add',title:'添加菜单'};
        $scope.saveMenu();
    }
    //{method:'PUT',url:''}

    $scope.editMenus = function(m){
        $scope.selectClassify=null;
        $scope.classifyChild=null;

        ActionTarget={method:'PUT',url:'content/'+m.ID,title:'修改菜单'};
        $scope.Menus=m;

        $scope.classify = {ContentID:$scope.Menus.ID};

        $("#editMenus").modal({centered: false,allowMultiple: true}).modal("show");

        $scope.listClassify();
    }
    $scope.deleteMenus = function(ID){
        $http.delete("content/"+ID,{transformRequest:angular.identity,headers:{"Content-Type":"application/json;charset=utf-8"}}).then(function(data){

            alert(data.data.Message);

            $scope.listMenus();

            $scope.Menus=null;

        });
    }
    $scope.listMenus();


    $scope.classify=null;

    $scope.ActionClassifyTarget={method:'POST',url:'content/sub_type',title:'添加分类'};

    $scope.deleteClassify = function(m){

        $http.delete("content/sub_type/"+m.ID,{transformRequest:angular.identity,headers:{"Content-Type":"application/json;charset=utf-8"}}).then(function(data){

            alert(data.data.Message);
            $scope.listClassify();

        });

    }
    $scope.editClassify = function(m){

        $scope.classify=m;
        $scope.ActionClassifyTarget={method:'PUT',url:'content/sub_type/'+m.ID,title:'修改分类'};
        //$scope.saveClassify();
    }
    $scope.saveClassify = function () {

        $http({
            method:$scope.ActionClassifyTarget.method,
            url:$scope.ActionClassifyTarget.url,
            data:JSON.stringify($scope.classify),
            transformRequest:angular.identity,
            headers:{"Content-Type":"application/json;charset=utf-8"}
        }).then(function(data){

            $scope.listClassify();
            $scope.classify.Name='';
            $scope.classify.ID=null;
            $scope.ActionClassifyTarget={method:'POST',url:'content/sub_type',title:'添加分类'};
            alert(data.data.Message);
        });


    }





    $scope.selectClassifyFunc = function(m){
        //$scope.ActionClassifyTarget={method:'PUT',url:'content_sub_type/'+m.ID,title:'修改分类'};
        $scope.selectClassify=m;
        $scope.listChildClassify($scope.selectClassify.ContentID,$scope.selectClassify.ID);
    }


    $scope.selectClassify=null;
    $scope.classifyChild=null;

    $scope.ActionClassifyChildTarget={method:'POST',url:'content/sub_type',title:'添加子分类'};

    //saveClassifyChild
    $scope.saveClassifyChild = function () {

        if(!$scope.selectClassify){
            alert("请选择父类");
            return
        }
        if(!$scope.Menus){
            alert("请菜单");
            return
        }


        $scope.classifyChild.ParentContentSubTypeID=$scope.selectClassify.ID;
        //{ContentID:$scope.Menus.ID};
        $scope.classifyChild.ContentID=$scope.Menus.ID;
        //$scope.classifyChild.MenusID=$scope.Menus.ID;
        //{MenusID:$scope.Menus.ID}

        $http({
            method:$scope.ActionClassifyChildTarget.method,
            url:$scope.ActionClassifyChildTarget.url,
            data:JSON.stringify($scope.classifyChild),
            transformRequest:angular.identity,
            headers:{"Content-Type":"application/json;charset=utf-8"}
        }).then(function(data){
            alert(data.data.Message);
            $scope.listChildClassify($scope.selectClassify.ContentID,$scope.selectClassify.ID);

            $scope.classifyChild.Name='';
            $scope.classifyChild.ID=null;
            $scope.ActionClassifyChildTarget={method:'POST',url:'content/sub_type',title:'添加分类'};
        });


    }
    $scope.deleteClassifyChild = function(m){

        $http.delete("content/sub_type/"+m.ID,{transformRequest:angular.identity,headers:{"Content-Type":"application/json;charset=utf-8"}}).then(function(data){

            alert(data.data.Message);
            $scope.listChildClassify($scope.selectClassify.ContentID,$scope.selectClassify.ID);

        });

    }
    $scope.editClassifyChild = function(m){
        $scope.classifyChild=m;
        $scope.ActionClassifyChildTarget={method:'PUT',url:'content/sub_type/'+m.ID,title:'修改分类'};
        //$scope.saveClassify();
    }



});


main.controller("store_situation_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    //store_journal/list
    $scope.StartTime =new Date();
    $scope.EndTime =new Date();

    $scope.tabIndex =1;
    $scope.selectTab= function (index) {
        $scope.tabIndex =index;
        if(table_local!=undefined){
            table_local.ajax.reload();
        }
    };

    var table_local;
    $timeout(function () {

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Detail"},
                {data:"StoreID",searchable:false},
                {data:"Type",searchable:false,visible:false},
                {data:"Amount",searchable:false,render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"Balance",searchable:false,render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "store_journal/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    d.Customs =[];
                    var st = $scope.StartTime;
                    st.setHours(0,0,0,0);
                    var et = $scope.EndTime;
                    et.setHours(24,0,0,0);
                    d.Customs.push({Name:"CreatedAt",Value:">='"+$filter("date")(st,"yyyy-MM-dd HH:mm:ss")+"'"});
                    d.Customs.push({Name:"CreatedAt",Value:"<'"+$filter("date")(et,"yyyy-MM-dd HH:mm:ss")+"'"});

                    if($scope.StoreID!=undefined&&$scope.StoreID!=""){
                        d.columns[3].search.value="'"+$scope.StoreID+"'";
                    }

                    d.columns[4].search.value="'"+$scope.tabIndex+"'";
                    return JSON.stringify(d);
                }
            }
        });



    });

});
main.controller("carditem_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.StartTime =new Date();
    $scope.EndTime =new Date();

    $scope.tabIndex ="OrdersGoods";
    $scope.selectTab= function (index) {
        $scope.tabIndex =index;

        if(table_local!=undefined){
            table_local.ajax.reload();
        }
    };
    $scope.submit = function () {
        if(table_local!=undefined){
            table_local.ajax.reload();
        }
    }

    var table_local;
    $timeout(function () {


        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"Type",searchable:false,render:function (data, type, row){

                    if(data=="OrdersGoods"){
                        return "商品";
                    }else if(data=="ScoreGoods"){
                        return "积分商品";
                    }else if(data=="Voucher"){
                        return "卡卷";
                    }else{
                        return "无";
                    }
                    }},
                {data:"ID"},
                {data:"UserID"},
                {data:"Data",searchable:false,render:function (data, type, row){
                        //console.log(type);
                        console.log(row.Type);

                        var Data = JSON.parse(data);

                        if(row.Type=="OrdersGoods"){
                            Data.Goods=JSON.parse(Data.Goods)
                            Data.Specification=JSON.parse(Data.Specification)
                            return Data.Goods.Title+"-"+Data.Specification.Label+"("+(Data.Specification.Num*Data.Specification.Weight/1000)+"Kg)";
                        }else if(row.Type=="ScoreGoods"){
                            return Data.Name;
                        }else if(row.Type=="Voucher"){
                            return Data.Name;
                        }else{
                            return "无";
                        }
                    }},
                {data:"Quantity",searchable:false},
                {data:"UseQuantity",searchable:false},
                {data:"ExpireTime",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }},
                {data:"PostType",searchable:false,render:function (data,type,row) {
                       if(data==1){
                           return "邮寄";
                       }else{
                           return "线下核销";
                       }
                    }},
                {data:"CreatedAt",searchable:false,render:function (data,type,row) {
                        return $filter("date")(data,"medium");
                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "carditem/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    //d.columns[0].search.value=$scope.tabIndex;
                    d.columns[0].search.value="'"+$scope.tabIndex+"'";
                    d.Customs =[];

                    var st = $scope.StartTime;
                    st.setHours(0,0,0,0);
                    //console.log(new Date());

                    var et = $scope.EndTime;
                    et.setHours(24,0,0,0);
                    //console.log(et.getFullYear(),et.getMonth(),et.getDate());

                    d.Customs.push({Name:"CreatedAt",Value:">='"+$filter("date")(st,"yyyy-MM-dd HH:mm:ss")+"'"});
                    d.Customs.push({Name:"CreatedAt",Value:"<'"+$filter("date")(et,"yyyy-MM-dd HH:mm:ss")+"'"});
                    return JSON.stringify(d);
                }
            }
        });



    });

});
main.controller("view_situation_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.StartTime =new Date();
    $scope.EndTime =new Date();
    $scope.situation ={};




    $scope.submit = function () {

        var form ={};
        form.StartTime = $scope.StartTime.getTime();
        form.EndTime = $scope.EndTime.getTime();

        $http.post("situation",$.param(form), {
            transformRequest: angular.identity,
            //headers: {"Content-Type": "application/x-www-form-urlencoded"}
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

            $scope.situation =data.data.Data;

        });

    }
    $scope.submit();


    //situation


});

main.controller("user_setup_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.rank =null;
    $scope.ConfigurationKey_ScoreConvertGrowValue =1;


    $scope.showAddRank = function(){
        //add_rank
        $("#add_rank").modal("show");

    }

    $scope.saveLeveConfiguration = function(){


        var total = $scope.ConfigurationKey_BrokerageLeve1+$scope.ConfigurationKey_BrokerageLeve2+$scope.ConfigurationKey_BrokerageLeve3+$scope.ConfigurationKey_BrokerageLeve4+$scope.ConfigurationKey_BrokerageLeve5+$scope.ConfigurationKey_BrokerageLeve6;
        if(total!=100){
            if(total!=0){
                alert("分佣比例不正确，比例总和为100或0");
                return
            }
        }


        $http.post("configuration/change",JSON.stringify({K:1201,V:$scope.ConfigurationKey_BrokerageLeve1.toString()}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            //alert(data.data.Message);

            $http.post("configuration/change",JSON.stringify({K:1202,V:$scope.ConfigurationKey_BrokerageLeve2.toString()}), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/json"}
            }).then(function (data, status, headers, config) {

                //alert(data.data.Message);

                $http.post("configuration/change",JSON.stringify({K:1203,V:$scope.ConfigurationKey_BrokerageLeve3.toString()}), {
                    transformRequest: angular.identity,
                    headers: {"Content-Type": "application/json"}
                }).then(function (data, status, headers, config) {

                    //alert(data.data.Message);

                    $http.post("configuration/change",JSON.stringify({K:1204,V:$scope.ConfigurationKey_BrokerageLeve4.toString()}), {
                        transformRequest: angular.identity,
                        headers: {"Content-Type": "application/json"}
                    }).then(function (data, status, headers, config) {

                        //alert(data.data.Message);

                        $http.post("configuration/change",JSON.stringify({K:1205,V:$scope.ConfigurationKey_BrokerageLeve5.toString()}), {
                            transformRequest: angular.identity,
                            headers: {"Content-Type": "application/json"}
                        }).then(function (data, status, headers, config) {

                            //alert(data.data.Message);

                            $http.post("configuration/change",JSON.stringify({K:1206,V:$scope.ConfigurationKey_BrokerageLeve6.toString()}), {
                                transformRequest: angular.identity,
                                headers: {"Content-Type": "application/json"}
                            }).then(function (data, status, headers, config) {

                                alert(data.data.Message);

                            });

                        });

                    });

                });

            });

        });


    }

    $scope.saveConfiguration = function(k,v){

        $http.post("configuration/change",JSON.stringify({K:k,V:v}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

        });
    }
    $scope.saveRank = function(){

        $http.post("rank/add",JSON.stringify($scope.rank), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            $("#add_rank").modal("hide");

            table_local.ajax.reload();
            $scope.rank =null;


        });

        //configuration/list

    }
    $scope.configurations={}
    $http.post("configuration/list",JSON.stringify([1100,1201,1202,1203,1204,1205,1206]), {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {

        var obj =data.data.Data;
        $scope.configurations=obj;
       // console.log(data.data.Data);
        $scope.ConfigurationKey_ScoreConvertGrowValue =parseInt(obj[1100]);

        $scope.ConfigurationKey_BrokerageLeve1=parseInt(obj[1201]);
        $scope.ConfigurationKey_BrokerageLeve2=parseInt(obj[1202]);
        $scope.ConfigurationKey_BrokerageLeve3=parseInt(obj[1203]);
        $scope.ConfigurationKey_BrokerageLeve4=parseInt(obj[1204]);
        $scope.ConfigurationKey_BrokerageLeve5=parseInt(obj[1205]);
        $scope.ConfigurationKey_BrokerageLeve6=parseInt(obj[1206]);
    });
    var table_local;
    $timeout(function () {


        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"GrowMaxValue"},
                {data:"Title"},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "initComplete":function (d) {

            },
            "ajax": {
                "url": "rank/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            //console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("rank/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });

    });

});

main.controller("order_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    $scope.tabIndex = parseInt(window.localStorage.getItem("TabIndex"));
    if(!$scope.tabIndex){
        $scope.tabIndex=0;
    }
    $scope.tabs=[
        {lable:"所有",status:""},
        {lable:"待付款",status:"Order"},
        {lable:"待发货",status:"Pay"},
        {lable:"待收货",status:"Deliver"},
        {lable:"申请退货退款",status:"Refund"},
        {lable:"已完成退货",status:"RefundOk"},
        {lable:"订单完成",status:"OrderOk"},
        {lable:"申请订单取消",status:"Cancel"},
        {lable:"已取消",status:"CancelOk"}
    ];
    $scope.currentTab =$scope.tabs[$scope.tabIndex];
    $scope.selectTab = function (index) {
        $scope.tabIndex =index;
        $scope.currentTab =$scope.tabs[$scope.tabIndex];

        window.localStorage.setItem("TabIndex",$scope.tabIndex);

        if(table_local){
            table_local.ajax.reload();
        }
    }
    var table_local;
    $timeout(function () {

        //UserID uint64, PostType int, Status

        table_local = $('#table_local').DataTable({
            //searching: false,
            "columns": [
                {data:"User.ID",orderable:false,render:function (data, type, row){return "";}},
                {data:"Orders.PostType",orderable:false,render:function (data, type, row){return "";}},
                {data:"Orders.Status",orderable:false,render:function (data, type, row){return "";}},
                {data:null,orderable:false,render:function (data, type, row){return "";}},
                {data:null,orderable:false,render:function (data, type, row){return "";}}

            ],
            createdRow:function ( row, data, index ) {
                //console.log(row);
                //console.log(data);
                //console.log(index);
                //$(row).hide();
            },
            drawCallback:function (settings) {

            },
            "rowCallback": function(row,data) {
                var rowsdfsdf = table_local.row(row);
                //console.log(row);



                console.log(row);

                var html =$('<div class="rowContent"></div>');


                var top =$('<div class=""></div>');

                if(data.Orders.PostType==1){
                    top =$('<div class="top post"></div>');
                }else if(data.Orders.PostType==2){
                    top =$('<div class="top store"></div>');
                }



                var info =$('<div class="info"></div>');
                info.text('订单#ID：'+data.Orders.ID);
                top.append(info);


                var info =$('<div class="info"></div>');
                info.html($filter("date")(data.Orders.CreatedAt,"medium"));
                top.append(info);


                var info =$('<div class="info"></div>');
                info.text(data.User.Name+"/"+data.User.Tel);
                top.append(info);

                var info =$('<div class="info"></div>');
                info.text('订单号：'+data.Orders.OrderNo);
                top.append(info);

                var info =$('<div class="info"></div>');
                info.text(data.Orders.IsPay==1?'支付成功':'未支付');
                top.append(info);

                //(data.Orders.IsPay==1?'支付':'未支付')

                html.append(top);


                var table = $('<table></table>');


                for(var i=0;i<data.OrdersGoodsList.length;i++){

                    var ordersGoods = data.OrdersGoodsList[i];

                    var Specification = JSON.parse(ordersGoods.Specification);
                    var Goods = JSON.parse(ordersGoods.Goods);
                    Goods.Images = JSON.parse(Goods.Images);

                    var tr = $('<tr data-index="'+i+'"></tr>');

                    var td = $('<td></td>');
                    var img = $("<img>");
                    img.attr("src",'/file/load?path='+Goods.Images[0]);
                    img.attr("width","100");
                    img.attr("height","100");
                    td.append(img);
                    tr.append(td);

                    var title =$('<td style="text-align: left;"></td>');

                    title.append('<div>'+Goods.Title+'</div>');
                    title.append('<div>规格：'+Specification.Label+'/'+(Specification.Num*Specification.Weight/1000)+'Kg</div>');
                    title.append('<div>'+('原价：'+Specification.CostPrice/100+'元，'+'市价：'+Specification.MarketPrice/100+'元，'+'分佣：'+Specification.Brokerage/100)+'元</div>');
                    tr.append(title);

                    var price =$('<td></td>');
                    price.append('<div style="color:#999;"><del>原价：'+(ordersGoods.CostPrice/100)+'元</del></div>');
                    price.append('<div>现价：'+(ordersGoods.SellPrice/100)+'元</div>');
                    tr.append(price);


                    var num =$('<td></td>');
                    num.append('<div>数量：'+(ordersGoods.Quantity)+'</div>');
                    tr.append(num);



                    var num =$('<td></td>');
                    num.append('<div><b>总金额：'+(ordersGoods.SellPrice*ordersGoods.Quantity/100)+'</b></div>');
                    tr.append(num);



                    if(i==0){
                        var num =$('<td class="operation" rowspan="99"></td>');
                        switch (data.Orders.Status){
                            case "Order":
                                //('<button class="ui blue mini button">修改支付金额</button>')
                                num.append('<div><button disabled class="ui mini button">等待支付</button><button class="ui blue PayMoney mini button">修改支付金额</button></div>');
                                break;
                            case "Pay":
                                if(data.Orders.PostType==1){
                                    num.append('<div><button class="ui red Deliver button">发货</button><button class="ui blue Cancel button">取消用户订单</button></div>');
                                }
                                break;
                            case "Refund":
                                num.append('<div><button disabled class="ui button">部分商品退款中</button></div>');
                                break;
                            case "Deliver":
                                num.append('<div><button disabled class="ui button">等待收货</button></div>');
                                break;
                            case "Cancel":
                                num.append('<div><button class="ui CancelOk blue button">处理取消申请</button></div>');
                                break;
                            case "CancelOk":
                                num.append('<div><button disabled class="ui button">取消成功</button></div>');
                                break;
                            case "OrderOk":
                                num.append('<div><button disabled class="ui button">订单完成</button></div>');
                                break;
                        }


                        if(data.Orders.PostType==1){

                            num.append('<div style="margin: 10px 0px;color:#666;">邮寄商品</div>');

                        }else if(data.Orders.PostType==2){
                            num.append('<div style="margin: 10px 0px;color:#666;">线下商品</div>');
                        }

                        ////是否支付，0=未支付，1=支付
                        //var ispay =$('<div><label>'+(data.Orders.IsPay==1?'支付':'未支付')+'</label></div>');
                        //num.append(ispay);

                        tr.append(num);



                        //var info =$('<div class="info"></div>').text("状态："+(data.Orders.Status));
                        //footer.append(info);
                    }






                    table.append(tr);




                    var tr = $('<tr data-index="'+i+'" class="tip"></tr>');
                    var num =$('<td colspan="5"></td>');

                    if(ordersGoods.Status=="OGAskRefund"){


                        var RefundInfo = JSON.parse(ordersGoods.RefundInfo);


                        var content=$('<div class="content"></div>');


                        var div = $('<div></div>');
                        div.text(RefundInfo.Reason);
                        content.append(div);

                        var div = $('<div></div>');

                        //包含货
                        if(RefundInfo.HasGoods){
                            div.append('<button class="ui blue RefundOk mini button">允许退货</button>');
                            div.append('<button class="ui RefundNo red mini button">拒绝申请</button>');
                        }else{
                            div.append('<button class="ui blue RefundOk mini button">允许退货</button>');
                            div.append('<button class="ui blue RefundComplete mini button">允许退款</button>');
                            div.append('<button class="ui RefundNo red mini button">拒绝申请</button>');
                        }


                        content.append(div);

                        num.append(content);


                    }else if(ordersGoods.Status=="OGRefundNo"){

                        var content=$('<div class="content"></div>');


                        var div = $('<div></div>');
                        content.append(div);

                        var div = $('<div>已经拒绝用户申请</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        content.append(div);


                        num.append(content);

                    }else if(ordersGoods.Status=="OGRefundOk"){


                        var content=$('<div class="content"></div>');

                        var div = $('<div></div>');
                        content.append(div);

                        var div = $('<div>已经同意用户退货申请，等待用户退货</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        content.append(div);

                        num.append(content);
                    }else if(ordersGoods.Status=="OGRefundInfo"){


                        var RefundInfo = JSON.parse(ordersGoods.RefundInfo);

                        var content=$('<div class="content"></div>');

                        var div = $('<div></div>');
                        div.append('<div>快递名称：'+RefundInfo.ShipName+'</div>');
                        div.append('<div>快递编号：'+RefundInfo.ShipNo+'</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        div.append(RefundInfo.ShipName);
                        content.append(div);




                        var div = $('<div><button class="ui red RefundComplete mini button">收到退货商品</button></div>');
                        content.append(div);

                        num.append(content);
                    }else if(ordersGoods.Status=="OGRefundComplete"){
                        var content=$('<div class="content"></div>');

                        var div = $('<div></div>');
                        content.append(div);

                        var div = $('<div>单品退货退款完成</div>');
                        content.append(div);

                        var div = $('<div></div>');
                        content.append(div);

                        num.append(content);
                    }

                    tr.append(num);
                    table.append(tr);

                }


                html.append(table);

                var  footer =$('<div class="footer"></div>');


                var info =$('<div class="info"></div>').text("商品金额："+(data.Orders.GoodsMoney/100)+"元");
                footer.append(info);

                var info =$('<div class="info"></div>').text("运费："+(data.Orders.ExpressMoney/100)+"元");
                footer.append(info);

                var info =$('<div class="info"></div>').text("优惠金额："+(data.Orders.DiscountMoney/100)+"元");
                footer.append(info);

                var info =$('<div style="color:blue;" class="info"></div>').text("总金额："+((data.Orders.GoodsMoney+data.Orders.ExpressMoney-data.Orders.DiscountMoney)/100)+"元");
                footer.append(info);

                var info =$('<div style="color:red;font-weight: bold;" class="info"></div>').text("支付金额："+(data.Orders.PayMoney/100)+"元");
                footer.append(info);


                var Address=JSON.parse(data.Orders.Address);
                var info =$('<div style="width: 250px;" class="info"></div>').text("邮寄地址："+(Address.Name+","+Address.Tel+","+Address.ProvinceName+Address.CityName+Address.CountyName+Address.Detail+","+Address.PostalCode));
                footer.append(info);




                html.append(footer);


                var info =$('<td colspan="99"></td>');
                info.append(html);

                $(row).empty().append(info);

                //rowsdfsdf.child(html).show();
            },
            "initComplete":function (d) {},
            "ajax": {
                "url": "order/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    d.columns[2].search.value=$scope.currentTab.status;
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.operation .PayMoney', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());


            $scope.$apply(function () {

                $scope.currentOrders=row.data();
                $("#ChangePayMoney").modal("show");

            });

        });
        $('#table_local').on('click','td.operation .Deliver', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());
            $scope.$apply(function () {

                $scope.currentOrders=row.data();
                $("#Deliver").modal("show");

            });

        });
        $('#table_local').on('click','td.operation .Cancel', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());
            $scope.$apply(function () {

                $scope.currentOrders=row.data();


                if(confirm("确定要取消用户的这个订单？")){
                    var form ={};
                    form.Action="Cancel";
                    form.ID=$scope.currentOrders.Orders.ID;
                    $http({
                        method:"PUT",
                        url:"order/change",
                        data:$.param(form),
                        transformRequest: angular.identity,
                        headers: {'Content-Type':'application/x-www-form-urlencoded'}
                    }).then(function (data, status, headers, config) {
                        alert(data.data.Message);
                        if(data.data.Success){
                            if(table_local){
                                table_local.ajax.reload();
                            }
                        }
                    });
                }

                //$("#CancelOk").modal("show");

            });

        });
        $('#table_local').on('click','td.operation .CancelOk', function () {
            var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);
            console.log(row.data());
            $scope.$apply(function () {

                $scope.currentOrders=row.data();
                $("#CancelOk").modal("show");

            });

        });
        $('#table_local').on('click','tr.tip .RefundOk', function () {
            var tr = $(this).closest('tr[role=row]');
            var OrdersGoodsIndex = tr.find(".rowContent").find(".tip").data("index");
            var row = table_local.row(tr);
            $scope.currentOrdersGoods=row.data().OrdersGoodsList[OrdersGoodsIndex];

            var form ={};
            form.Action="RefundOk";
            form.OrdersGoodsID= $scope.currentOrdersGoods.ID;
            $http({
                method:"PUT",
                url:"order/change",
                data:$.param(form),
                transformRequest: angular.identity,
                headers: {'Content-Type':'application/x-www-form-urlencoded'}
            }).then(function (data, status, headers, config) {
                alert(data.data.Message);
                if(table_local){
                    table_local.ajax.reload();
                }
            });

        });
        $('#table_local').on('click','tr.tip .RefundNo', function () {

            var tr = $(this).closest('tr[role=row]');
            var OrdersGoodsIndex = tr.find(".rowContent").find(".tip").data("index");
            var row = table_local.row(tr);
            $scope.currentOrdersGoods=row.data().OrdersGoodsList[OrdersGoodsIndex];




            var form ={};
            form.Action="RefundNo";
            form.OrdersGoodsID=$scope.currentOrdersGoods.ID;
            $http({
                method:"PUT",
                url:"order/change",
                data:$.param(form),
                transformRequest: angular.identity,
                headers: {'Content-Type':'application/x-www-form-urlencoded'}
            }).then(function (data, status, headers, config) {
                alert(data.data.Message);
                if(table_local){
                    table_local.ajax.reload();
                }
            });

        });
        $('#table_local').on('click','tr.tip .RefundComplete', function () {

            var tr = $(this).closest('tr[role=row]');
            var OrdersGoodsIndex = tr.find(".rowContent").find(".tip").data("index");
            var row = table_local.row(tr);

            //RefundComplete

            $scope.$apply(function () {
                $scope.currentOrders=row.data().Orders;
                $scope.currentOrdersGoods=row.data().OrdersGoodsList[OrdersGoodsIndex];
                $("#RefundComplete").modal({closable:false}).modal("show");

            });


            /*var tr = $(this).closest('tr[role=row]');
            var row = table_local.row(tr);

            if(confirm("退款将该单品金额退回给用户，如有参加满减活动，则按比例扣除金额。")){


                var OrdersGoodsID = tr.find(".rowContent").find(".tip").data("id");
                var RefundType = tr.find(".rowContent").find(".tip").find(".RefundComplete").data("refundtype");

                var form ={};
                form.Action="RefundComplete";
                form.OrdersGoodsID=OrdersGoodsID;
                form.RefundType=parseInt(RefundType);
                $http({
                    method:"PUT",
                    url:"order/change",
                    data:$.param(form),
                    transformRequest: angular.identity,
                    headers: {'Content-Type':'application/x-www-form-urlencoded'}
                }).then(function (data, status, headers, config) {
                    alert(data.data.Message);
                    if(table_local){
                        table_local.ajax.reload();
                    }
                });

            }*/



        });

    });


    $scope.currentOrders={};
    $scope.PayMoney=-1;
    $scope.ChangePayMoney = function () {

        if($scope.PayMoney<0){
            alert("请输入正确的金额");
            return;
        }


        var form ={};
        form.Action="PayMoney";
        form.PayMoney=$scope.PayMoney;
        form.ID=$scope.currentOrders.Orders.ID;
        $http({
            method:"PUT",
            url:"order/change",
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Success){
                $("#ChangePayMoney").modal("hide");
                if(table_local){
                    table_local.ajax.reload();
                }
            }
        });

    }
    $scope.ShipName=null;
    $scope.ShipNo=null;
    $scope.DeliverSubmit = function () {

        if($scope.ShipName==""){

            return;
        }
        if($scope.ShipNo==""){

            return;
        }


        var form ={};
        form.Action="Deliver";
        form.ShipName=$scope.ShipName;
        form.ShipNo=$scope.ShipNo;
        form.ID=$scope.currentOrders.Orders.ID;
        $http({
            method:"PUT",
            url:"order/change",
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Success){
                $("#Deliver").modal("hide");
                if(table_local){
                    table_local.ajax.reload();
                }
            }
        });

    }

    $scope.RefundType = 0;
    $scope.RefundCompleteSubmit = function(){



        if(confirm("退款将该单品金额退回给用户，如有参加满减活动，则按比例扣除金额。")){

            var form ={};
            form.Action="RefundComplete";
            form.OrdersGoodsID=$scope.currentOrdersGoods.ID;
            form.RefundType=parseInt($scope.RefundType);
            $http({
                method:"PUT",
                url:"order/change",
                data:$.param(form),
                transformRequest: angular.identity,
                headers: {'Content-Type':'application/x-www-form-urlencoded'}
            }).then(function (data, status, headers, config) {
                alert(data.data.Message);

                if(data.data.Success){

                    $("#RefundComplete").modal("hide");
                    $scope.currentOrders=null;
                    $scope.currentOrdersGoods=null;
                    if(table_local){
                        table_local.ajax.reload();
                    }
                }

            });

        }

    }

    $scope.RefundType = 0;
    $scope.CancelOkSubmit = function () {


        var form ={};
        form.Action="CancelOk";
        form.ID=$scope.currentOrders.Orders.ID;
        form.RefundType=$scope.RefundType;
        $http({
            method:"PUT",
            url:"order/change",
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Success){
                $("#CancelOk").modal("hide");
                if(table_local){
                    table_local.ajax.reload();
                }
            }
        });

    }



})
main.controller("add_express_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {


    $scope.TypeN={Type:"N",Areas:[],N:0};
    $scope.TypeM={Type:"M",Areas:[],M:0};
    $scope.TypeNM={Type:"NM",Areas:[],N:0,M:0};

    $scope.Template={Name:'',Type:'ITEM',Drawee:"BUYERS"};
    $scope.FreeItems=[];
    $scope.FreeItem={Areas:[],Type:'N'};
    $scope.defaultFreeItem=angular.copy($scope.FreeItem);
    $scope.deleteFreeItem = function(index){
        if(confirm("确定删除？")){
            $scope.FreeItems.splice(index,1);
        }
    }
    $scope.addFreeItem = function(){

        var Type=$scope.FreeItem.Type;
        if(Type=="N"){
            if($scope.FreeItem.N<=0){
                return
            }
        }
        if(Type=="M"){
            if($scope.FreeItem.M<=0){
                return
            }
        }
        if(Type=="NM"){
            if($scope.FreeItem.N<=0){
                return
            }
            if($scope.FreeItem.M<=0){
                return
            }
        }
        if($scope.FreeItem.Areas.length<=0){
            return
        }

        $scope.FreeItems.push($scope.FreeItem);

        $scope.FreeItem=angular.copy($scope.defaultFreeItem);

    }

    $scope.saveExpress = function(){


        //express_template/save

        for(var i=0;i<$scope.FreeItems.length;i++){
            var item =$scope.FreeItems[i];

            if(item.Type=="N"){
                if(item.N<=0){
                    alert("数据不完整");
                    return
                }
            }
            if(item.Type=="M"){
                if(item.M<=0){
                    alert("数据不完整");
                    return
                }
            }
            if(item.Type=="NM"){
                if(item.N<=0){
                    alert("数据不完整");
                    return
                }
                if(item.M<=0){
                    alert("数据不完整");
                    return
                }
            }
            if(item.Areas.length<=0){
                alert("数据不完整");
                return
            }
        }



        var dfd=$scope.jcsj($scope.items.Default);
        if(dfd==false){
            alert("数据不完整");
            return
        }
        for(var i=0;i<$scope.items.Items.length;i++){
            var item  = $scope.items.Items[i];
            if(item.Areas.length<=0){
                alert("数据不完整");
                return false
            }
            var ii = $scope.jcsj(item);
            if(ii==false){
                //alert("数据不完整");
                return
            }
        }

        var Template = {};
        Template.ID = $scope.Template.ID;
        Template.Name = $scope.Template.Name;
        Template.Type = $scope.Template.Type;
        Template.Drawee = $scope.Template.Drawee;
        Template.Template =JSON.stringify($scope.items);
        Template.Free =JSON.stringify($scope.FreeItems);

        $http.post("express_template/save",JSON.stringify(Template), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);
            if(data.data.Success==true){
                window.history.back();
            }

        });
    }

    $scope.jcsj = function(item){

        if(item.N<=0){
            alert("数据不完整");
            return false
        }
        if(item.M<=0){
            alert("数据不完整");
            return false
        }
        if(item.AN<=0){
            alert("数据不完整");
            return false
        }
        if(item.ANM<=0){
            alert("数据不完整");
            return false
        }

        return true;
    }
    $scope.dq =["上海","江苏省","浙江省","安徽省","江西省","北京","天津","山西省","山东省","河北省","内蒙古自治区","湖南省","湖北省","河南省","广东省","广西壮族自治区","福建省","海南省","辽宁省","吉林省","黑龙江省","陕西省","新疆维吾尔自治区","甘肃省","宁夏回族自治区","青海省","重庆","云南省","贵州省","西藏自治区","四川省"];


    //ng-checked="Template.FreeN.Enable"
    $scope.items ={Default:{Areas:[],N:0,M:0,AN:0,ANM:0},Items:[]};
    $scope.copyDefault = angular.copy($scope.items.Default);

    $scope.deleteItem = function(index){
        $scope.items.Items.splice(index,1);
    }
    $scope.addItem = function(){
        $scope.items.Items.push(angular.copy($scope.copyDefault));
    }

    $scope.currentItemIndex=-1;
    $scope.selectArea = function(index){

        //alert($scope.dq[index]);
        var area = $scope.dq[index];
        var areaIndex = $scope.items.Items[$scope.currentItemIndex].Areas.indexOf(area);
        if(areaIndex!=-1){
            $scope.items.Items[$scope.currentItemIndex].Areas.splice(areaIndex,1);
        }else{
            $scope.items.Items[$scope.currentItemIndex].Areas.push(area);
        }
        console.log($scope.items.Items);
    }
    $scope.AreaIndexList = [];
    $scope.addArea = function(index){
        $scope.currentItemIndex=index;
        $scope.AreaIndexList=[];

        for(var i=0;i<$scope.items.Items.length;i++){
            if(i!=$scope.currentItemIndex){
                var Areas = $scope.items.Items[i].Areas;
                for(var o=0;o<Areas.length;o++){
                    var area = Areas[o];
                    var areaIndex = $scope.AreaIndexList.indexOf(area);
                    if(areaIndex==-1){
                        $scope.AreaIndexList.push(area);
                    }

                }

            }

        }
        console.log($scope.AreaIndexList);

        $("#area_item").modal("show");
    }


    $scope.AreaTjIndexList=[];
    var AreaTjIndex =-1;
    $scope.addAreaJT = function(index){
        AreaTjIndex =index;
        //$scope.FreeItems=[];
        $scope.AreaTjIndexList=[];

        for(var i=0;i<$scope.FreeItems.length;i++){

            if(i!=index){
                $scope.AreaTjIndexList=$scope.AreaTjIndexList.concat($scope.FreeItems[i].Areas)
            }
        }

        console.log($scope.AreaTjIndexList);
        $("#area_tj").modal("show");
    }
    //selectAreaTJ

    $scope.selectAreaTJ = function(areaTxt){

        if(AreaTjIndex==-1){
            var areaIndex =$scope.FreeItem.Areas.indexOf(areaTxt);
            if(areaIndex==-1){
                $scope.FreeItem.Areas.push(areaTxt);
            }else{
                $scope.FreeItem.Areas.splice(areaIndex,1);
            }

        }else{
            var areaIndex =$scope.FreeItems[AreaTjIndex].Areas.indexOf(areaTxt);
            if(areaIndex==-1){
                $scope.FreeItems[AreaTjIndex].Areas.push(areaTxt);
            }else{
                $scope.FreeItems[AreaTjIndex].Areas.splice(areaIndex,1);
            }
        }

        //alert($scope.dq[index]);
        //var area = $scope.dq[index];
        /*var areaIndex = $scope.Template[TargetFree].Areas.indexOf(areaTxt);
        if(areaIndex!=-1){
            $scope.Template[TargetFree].Areas.splice(areaIndex,1);
        }else{
            $scope.Template[TargetFree].Areas.push(areaTxt);
        }
        console.log($scope.Template);*/
    }


    //alert($routeParams.ID);



    if($routeParams.ID!=undefined){

        $http.get("express_template/"+$routeParams.ID,JSON.stringify({}), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

           var et = data.data.Data;



            var Template = {};
            Template.ID = et.ID;
            Template.Name = et.Name;
            Template.Type = et.Type;
            Template.Drawee = et.Drawee;

            $scope.FreeItems =JSON.parse(et.Free);

            $scope.items =JSON.parse(et.Template);

            $scope.Template = Template;


        });
    }

    //express_template/:ID

    $timeout(function () {
        //$('.ui.radio.checkbox').checkbox();
        //$('.ui.checkbox').checkbox();
        //$(".ui.modal").modal("show");
    });
});
main.controller("express_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {


    var table_local;
    $timeout(function () {

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Drawee"},
                {data:"Type"},
                {data:"Free",orderable:false,render:function (data, type, row) {

                        var m = {};
                        try {
                            m = JSON.parse(data)
                        }catch (e) {

                        }
                        return m.length>0?'是':'否';

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui edit blue mini button">修改/查看</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "initComplete":function (d) {
                var info = table_local.page.info();
                var dataRows = info.recordsTotal;
                if(dataRows>0){
                    $("#add_express_btn").hide();
                }else{
                    $("#add_express_btn").show();
                }
            },
            "ajax": {
                "url": "express_template/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("express_template/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });
        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());

            //$scope.Admin=row.data();

            window.location.href="#!/add_express?ID="+row.data().ID;

            //$scope.showAdminModal({method:'PUT',url:'admin/'+$scope.Admin.ID,title:'修改密码'});

        });

    });

});
main.controller("admin_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.TargetAction={method:"",url:"",title:""};

    $scope.Admin=null;

    var table_local;

    $scope.saveAdmin = function () {


        $http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data: JSON.stringify($scope.Admin),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            table_local.ajax.reload();

            $("#adminModal").modal("hide");

            $scope.Admin=null;
            $scope.PassWord=null;

        });




    }
    $scope.showAdminModal = function (targetAction) {

        $timeout(function () {
            $scope.TargetAction=targetAction;
            $("#adminModal").modal("show");
        });
    }

    $timeout(function () {
        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Account",render:function (data) {
                        if(data==LoginAccount){
                            return data+"【自己】";
                        }else{
                            return data;
                        }
                    }},
                {data:"LastLoginAt",render:function (data) {

                    return $filter("date")(data,"medium");

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        if(LoginAccount=="admin"){
                            return '<button class="ui edit blue mini button">修改密码</button>'+
                                '<button class="ui authority teal mini button">权限管理</button>'+
                                '<button class="ui delete red mini button">删除</button>';
                        }else{
                            if(data.Account==LoginAccount){
                                return '<button class="ui edit blue mini button">修改密码</button>';
                            }else{
                                return '';
                            }

                        }


                    }}
            ],
            "ajax": {
                "url": "admin/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });
        $('#table_local').on('click','td.opera .authority', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());

            //$scope.showAdminModal({method:'POST',url:'admin',title:'添加管理员'});

            //authorityModal

            $http.get("admin/"+row.data().ID,{}).then(function (data) {

                var authoritys =[];
                $scope.Admin=data.data.Data;
                authoritys = JSON.parse($scope.Admin.Authority);

                $('#authorityModal .ui.toggle.checkbox').checkbox("set unchecked");

                for(var i=0;i<authoritys.length;i++){

                    var name = authoritys[i];
                    $('#authorityModal .ui.toggle.checkbox input[name='+name+']').parent().checkbox("set checked");
                }

                // console.log($('#authorityModal .ui.toggle.checkbox input[name='+key+']'))

                $("#authorityModal").modal({centered:false,onApprove:function () {


                        $http({
                            method:"PUT",
                            url:'admin/authority/'+$scope.Admin.ID,
                            data: JSON.stringify({Authority:JSON.stringify(authoritys)}),
                            transformRequest: angular.identity,
                            headers: {"Content-Type": "application/json"}
                        }).then(function (data, status, headers, config) {

                            alert(data.data.Message);

                            $("#authorityModal").modal("hide");

                            $scope.Admin=null;


                        });


                        return false;

                    }}).modal('setting', 'closable', false).modal("show");



                $('#authorityModal .ui.toggle.checkbox').checkbox({
                    onChecked: function() {
                        //console.log($(this).data("value"));
                        //console.log(eval("("+$(this).data("value")+")"));
                        //authoritys[$(this).attr("name")]=eval("("+$(this).data("value")+")");
                        authoritys.push($(this).attr("name"));
                    },
                    onUnchecked: function() {
                        //console.log($(this).attr("name"));
                        //delete authoritys[$(this).attr("name")];
                        var name = $(this).attr("name");
                        var index = authoritys.indexOf(name);
                        authoritys.splice(index,1);

                        console.log(authoritys);
                    },
                });





            })

        });
        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("admin/"+row.data().ID,{}).then(function (data) {
                    alert(data.data.Message);
                    table_local.ajax.reload();

                })
            }
        });
        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());

            $scope.Admin=row.data();

            $scope.showAdminModal({method:'PUT',url:'admin/'+$scope.Admin.ID,title:'修改密码'});

        });
    });

});
main.controller("store_stock_manager_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.Store={};
    $scope.Store.ID=$routeParams.ID;
    if($scope.Store.ID==undefined){

        alert("无没有门店信息")
        window.history.back();
    }


    $scope.Specifications=[];

    $http({
        method:"GET",
        url:"store/"+$scope.Store.ID,
        data:{},
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        var Store = data.data.Data;
        $scope.Images=JSON.parse(Store.Images);
        $scope.Pictures=JSON.parse(Store.Pictures);
        $scope.Store=Store;

    });


    $scope.TargetAction={method:"",url:"",title:""};
    $scope.cancelStoreStock = function(){
        $scope.StoreStock={};
        $scope.TargetAction={method:"POST",url:"store/stock",title:"添加产品规格数量"}

    }
    $scope.AddStoreStockStock=0;
    $scope.saveStoreStock = function(){

        if($scope.SelectGoods==null){
            alert("请选择产品");
            return
        }

        if($scope.StoreStock.SpecificationID==undefined){
            alert("请选择产品规格");
            return
        }


        $scope.StoreStock.StoreID =parseInt($routeParams.ID);
        $scope.StoreStock.GoodsID=$scope.SelectGoods.ID;

        var form ={};
        form.StoreID=parseInt($routeParams.ID);
        form.GoodsID=$scope.SelectGoods.ID;
        form.ID=$scope.StoreStock.ID;
        form.SpecificationID=$scope.StoreStock.SpecificationID;
        form.AddStoreStockStock=$scope.AddStoreStockStock;

        $http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            $scope.StoreStock=null;

            alert(data.data.Message);

            $scope.cancelStoreStock();
            $scope.AddStoreStockStock=0;

            table_local_goods.ajax.reload(null,false);
            table_local_stock.ajax.reload(null,false);
            table_store_stock.ajax.reload(null,false);

        });


        //$("#add_store_stock").modal("hide");

        /*$http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            $scope.cancelStoreStock();

            table_local_goods.ajax.reload(null,false);
            table_local_stock.ajax.reload(null,false);
            table_store_stock.ajax.reload(null,false);

        });

        //$scope.SelectGoods=null;
        $scope.StoreStock=null;*/







    }

    $scope.StoreStockModal = function(targetAction){
        $scope.TargetAction = targetAction;

        $("#add_store_stock").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("show");

    }
    $scope.addGoods = function () {

        $("#list_goods").modal({centered:true,allowMultiple: false}).modal('setting', 'closable', false).modal("show");
        //$scope.StoreStockModal();
        //$scope.ListGoodsSpecification(2008);
    }

    $scope.SpecificationsDisable={};
    $scope.StoreGoodsExist={};
    $scope.ListGoodsSpecification = function (GoodsID) {


        $http.get("goods?action=get_goods",{params:{ID:GoodsID}}).then(function (data) {

            //alert(data.data.Message);
            //$scope.StoreStock=row.data();
            $scope.SelectGoods=data.data.Data.Goods;
            $scope.Specifications=data.data.Data.Specifications;
            // $scope.StoreStockModal({method:"PUT",url:"store/stock/"+$scope.StoreStock.ID,title:"修改门店库存"});



            if(table_store_stock!=null){

                table_store_stock.ajax.url("store/stock/list/"+$scope.Store.ID+"/"+GoodsID).load(null,false);
                return
            }
            table_store_stock = $('#table_store_stock').DataTable({
                searching:false,
                "createdRow": function( row, data, dataIndex ) {
                    //console.log(row,data,dataIndex);
                    var SpecificationsDisable = $scope.SpecificationsDisable;
                    SpecificationsDisable[data.StoreStock.SpecificationID]=true;
                    $scope.SpecificationsDisable=SpecificationsDisable;
                },
                "columns": [
                    {data:"StoreStock.ID"},
                    {data:"Goods.Title"},
                    {data:"Specification.Label"},
                    {data:"StoreStock.Stock",render:function (data, type, row) {
                            //console.log(row.StoreStock.Stock-row.StoreStock.UseStock)
                            return row.StoreStock.Stock-row.StoreStock.UseStock;

                        }},
                    {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                            return '<button class="ui edit blue mini button">编辑</button><button class="ui delete red mini button">删除</button>';

                        }}
                ],
                "ajax": {
                    "url": "store/stock/list/"+$scope.Store.ID+"/"+GoodsID,
                    "type": "POST",
                    "contentType": "application/json",
                    "data": function ( d ) {
                        //d.columns[1].search.value=$scope.Store.ID.toString();
                        return JSON.stringify(d);
                    }
                }
            });
            $('#table_store_stock').on('click','td.opera .edit', function () {
                var tr = $(this).closest('tr');
                var row = table_store_stock.row(tr);
                //console.log(row.data());

                $scope.StoreStock=row.data().StoreStock;
                //$scope.selectGoods=null;

                ///$scope.TargetAction={method:"POST",url:"store_stock",title:"产品规格数量"}


                //$scope.TargetAction={method:"PUT",url:"store_stock/"+row.data().StoreStock.ID,title:"修改产品规格数量"}
                $scope.TargetAction={method:"PUT",url:"store/stock",title:"修改产品规格数量"}
                $scope.ListGoodsSpecification(row.data().StoreStock.GoodsID);

            });

            $('#table_store_stock').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table_store_stock.row(tr);

                if(confirm("确定删除？")){
                    $scope.SpecificationsDisable=[];
                    $http.delete("store/stock/"+row.data().StoreStock.ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table_store_stock.ajax.reload(null,false);

                    })
                }
            });

        })
        //modal('attach events', '#add_store_stock .actions .button').
        $("#add_store_stock").modal({detachable:true,centered:true,allowMultiple: false}).modal('setting', 'closable', false).modal("show");
        //$("#list_goods").modal({centered:true,allowMultiple: false}).modal('setting', 'closable', false).modal("show");
        //$scope.StoreStockModal();
    }

    var table_local_goods;
    var table_local_stock;
    var table_store_stock;

    $timeout(function () {


        table_local_goods = $('#table_local_goods').DataTable({
            fixedColumns: true,
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        if($scope.StoreGoodsExist[data.ID]){
                            return '<button disabled class="ui blue mini button">已选</button>';
                        }else {
                            return '<button class="ui select blue mini button">选择</button>';
                        }
                    }}
            ],
            "ajax": {
                "url": "goods?action=list_goods",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local_goods').on('click','td.opera .select', function () {
            var tr = $(this).closest('tr');
            var row = table_local_goods.row(tr);

            $timeout(function () {
                //$("#list_goods").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("hide");

                $scope.SelectGoods=row.data();
                $scope.StoreStock=null;

                $scope.TargetAction={method:"POST",url:"store/stock",title:"添加产品规格数量"}
                //$scope.StoreStockModal({method:"POST",url:"store/stock",title:"添加门店库存"});

                //$scope.TargetAction={method:"POST",url:"store_stock",title:"产品规格数量"}

                //$scope.TargetAction={method:"POST",url:"store_stock",title:"产品规格数量"}
                $scope.ListGoodsSpecification($scope.SelectGoods.ID);

            });
        });


        table_local_stock = $('#table_local_stock').DataTable({
            searching:false,
            "columns": [
                {data:"GoodsID"},
                {data:"StoreID",visible:false},
                {data:"Title"},
                {data:"Total"},
                {data:"Stock",render:function (data,type,row) {
                        //console.log(row);
                        //row.Stock-row.UseStock
                        return row.Stock-row.UseStock;

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button class="ui edit blue mini button">编辑</button>';//<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "drawCallback": function(settings) {

                //store_stock/able/goods/:StoreID

                $http.get("store/stock/exist/goods/"+$scope.Store.ID).then(function (data) {
                    console.log(data.data.Data);
                    var list = data.data.Data;
                    var StoreGoodsExist = {};
                    for(var i=0;i<list.length;i++){
                        StoreGoodsExist[list[i].GoodsID] = true;
                    }
                    $scope.StoreGoodsExist=StoreGoodsExist;

                    table_local_goods.draw(false);
                })
            },
            "ajax": {
                "url": "store/stock/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    d.columns[1].search.value=$scope.Store.ID.toString();
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local_stock').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table_local_stock.row(tr);
            //console.log(row.data());

            //$scope.StoreStock=null;
            //$scope.selectGoods=null;

            $scope.TargetAction={method:"POST",url:"store/stock",title:"添加产品规格数量"}
            //$scope.TargetAction={method:"PUT",url:"store_stock/"+$scope.StoreStock.ID,title:"修改产品规格数量"}
            $scope.ListGoodsSpecification(row.data().GoodsID);

        });


    });

})
main.controller("store_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.selectGoods=null;
    $scope.StoreStock=null;
    $scope.Store=null;

    var table_local_goods;
    var table_local_stock;
    var table_local;


    $scope.TargetAction={method:"",url:"",title:""};

    $timeout(function () {

        table_local = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Phone"},
                {data:"Address",render:function (data, type, row) {
                    var address =JSON.parse(data);
                    //{"ProvinceName":"福建省","CityName":"三明市","CountyName":"梅列区",
                        // "Detail":"列东街道东新二路45号天元列东饭店",
                        // "PostalCode":"350402","Name":"fsdfdsfsd","Tel":"13809549424"}

                    return address.ProvinceName+address.CityName+address.CountyName+address.Detail;

                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<a href="#!/add_store?ID='+data.ID+'" class="ui edit blue mini button">编辑</a>'+
                            '<button class="ui add_goods_stock teal mini button">库存管理</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "ajax": {
                "url": "store/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#table_local').on('click','td.opera .delete', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            console.log(row.data());
            if(confirm("确定删除？")){
                $http.delete("store/"+row.data().ID,{params:{}}).then(function (data) {

                    alert(data.data.Message);

                    table_local.ajax.reload();

                })
            }
        });

        $('#table_local').on('click','td.opera .add_goods_stock', function () {
            var tr = $(this).closest('tr');
            var row = table_local.row(tr);
            //console.log(row.data());
            window.location.href="#!/store_stock_manager?ID="+row.data().ID;
            //$scope.Store=row.data();
            //$('#table_local_stock').DataTable().column(1).search($scope.Store.ID).draw();
            //$("#add_goods_stock").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("show");
            //$("#add_store_stock").modal({centered:true,allowMultiple: true}).modal('setting', 'closable', false).modal("show");

        });


    });

});
main.controller("add_store_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.Images=[];
    $scope.Pictures=[];

    $scope.Store ={ID:$routeParams.ID};


    $scope.TargetAction={method:"POST",url:"store/add",title:"添加门店"};

    if($scope.Store.ID!=undefined){

        $scope.TargetAction={method:"PUT",url:"store/"+$scope.Store.ID,title:"修改门店"};

        $http({
            method:"GET",
            url:"store/"+$scope.Store.ID,
            data:{},
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {
            var Store = data.data.Data;
            $scope.Images=JSON.parse(Store.Images);
            $scope.Pictures=JSON.parse(Store.Pictures);

            try {
                $scope.address=JSON.parse(Store.Address);
            }catch (e) {
                $scope.address={};
            }

            $scope.Store=Store;

        });


    }






    $scope.address ={};
    $scope.showMapModal =function(){


        var currentPositionResult;

        $("#selectMap").modal({centered: false,onApprove:function (e) {

                if(currentPositionResult==undefined){
                    alert("请选择坐标");
                    return false
                }

                $timeout(function () {
                    console.log("currentPositionResult",currentPositionResult)
                   // $scope.Store.Latitude=currentPositionResult.lat;
                    //$scope.Store.Longitude=currentPositionResult.lng;
                    $scope.Store.Latitude=currentPositionResult.lat;
                    $scope.Store.Longitude=currentPositionResult.lng;

                });

            }}).modal("show");




        AMapUI.loadUI(['misc/PositionPicker'], function(PositionPicker) {
            var map = new AMap.Map('container', {
                zoom: 16,
                scrollWheel: false
            })
            var positionPicker = new PositionPicker({
                mode: 'dragMarker',
                map: map
            });


            var infoWindow = new AMap.InfoWindow({
                autoMove: true,
                offset: {x: 0, y: -30}
            });



            if($scope.Store.Longitude && $scope.Store.Latitude){
                positionPicker.start(new AMap.LngLat($scope.Store.Longitude,$scope.Store.Latitude));
                map.panTo(new AMap.LngLat($scope.Store.Longitude,$scope.Store.Latitude));
            }


            //map.panTo(currentPositionResult);
            positionPicker.on('success', function(positionResult) {
                currentPositionResult=positionResult.position;

                console.log(positionResult);

                //infoWindow.setContent(createContent(poiArr[0]));
                //infoWindow.setContent($scope.Store.Address);
                infoWindow.open(map,currentPositionResult);

                $scope.address.ProvinceName=positionResult.regeocode.addressComponent.province;
                $scope.address.CityName=positionResult.regeocode.addressComponent.city;
                $scope.address.CountyName=positionResult.regeocode.addressComponent.district;
                $scope.address.Detail=positionResult.regeocode.addressComponent.township+positionResult.regeocode.addressComponent.street+positionResult.regeocode.addressComponent.streetNumber;
                if(positionResult.regeocode.pois.length>0){
                    $scope.address.Detail=$scope.address.Detail+positionResult.regeocode.pois[0].name;
                }

                $scope.address.PostalCode=positionResult.regeocode.addressComponent.adcode;


                infoWindow.setContent($scope.address.ProvinceName+$scope.address.CityName+$scope.address.CountyName+$scope.address.Detail);


                document.getElementById('lnglat').innerHTML = positionResult.position;
                document.getElementById('address').innerHTML = positionResult.address;
                document.getElementById('nearestJunction').innerHTML = positionResult.nearestJunction;
                document.getElementById('nearestRoad').innerHTML = positionResult.nearestRoad;
                document.getElementById('nearestPOI').innerHTML = positionResult.nearestPOI;
            });
            positionPicker.on('fail', function(positionResult) {
                document.getElementById('lnglat').innerHTML = ' ';
                document.getElementById('address').innerHTML = ' ';
                document.getElementById('nearestJunction').innerHTML = ' ';
                document.getElementById('nearestRoad').innerHTML = ' ';
                document.getElementById('nearestPOI').innerHTML = ' ';
            });



            var startButton = document.getElementById('start');
            var stopButton = document.getElementById('stop');
            var dragMapMode = document.getElementsByName('mode')[0];
            var dragMarkerMode = document.getElementsByName('mode')[1];

            //serachValue   serachBtn

            var serachValue = document.getElementById('serachValue');
            var serachBtn = document.getElementById('serachBtn');





            function PlaceSearch() {
                var serachTxt = $(serachValue).val();
                if(serachTxt==""){
                    alert("请输入地点名称");
                    return
                }
                //console.log($(serachValue).val());

                AMap.plugin('AMap.PlaceSearch', function(){
                //AMap.service(["AMap.PlaceSearch"], function() {
                    var placeSearch = new AMap.PlaceSearch({ //构造地点查询类
                        pageSize: 1,
                        pageIndex: 1,
                        //city: "010", //城市
                        //map: map,
                        //panel: "panel"
                    });
                    //关键字查询
                    placeSearch.search(serachTxt,function(status,result){
                        if(status=="complete"){
                            if(result.poiList.pois.length>0){
                                positionPicker.start(result.poiList.pois[0].location);
                                currentPositionResult=result.poiList.pois[0].location;
                            }else {
                                currentPositionResult=null;
                            }
                        }else{
                            alert(status);
                        }


                    });
                });
            }


            serachBtn.addEventListener("click",PlaceSearch)
            serachValue.addEventListener("keypress",PlaceSearch)
            startButton.addEventListener("click",function () {
                if(currentPositionResult){

                    map.panTo(currentPositionResult);
                }
            })




            positionPicker.start();
            map.panBy(0, 1);

            map.addControl(new AMap.ToolBar({
                liteStyle: true
            }))
        });


    }

    $scope.deleteArr = function(arr,index){
        if(confirm("确认删除这项内容？")){
            arr.splice(index,1);
        }
    }



    $scope.save = function(){

        $scope.Store.Images = JSON.stringify($scope.Images);
        $scope.Store.Pictures = JSON.stringify($scope.Pictures);

        if($scope.Store.Latitude==undefined||$scope.Store.Latitude==""||$scope.Store.Longitude==undefined||$scope.Store.Longitude==""){

            alert("请选择坐标地址");
            return
        }


        if($scope.address.Detail==undefined||$scope.address.Detail==""){

            alert("请选择填写地址");
            return
        }

        $scope.address.Name=$scope.Store.Name;
        $scope.address.Tel=$scope.Store.ServicePhone;

        $scope.Store.Address=JSON.stringify($scope.address);


        $http({
            method: $scope.TargetAction.method,
            url: $scope.TargetAction.url,
            data: JSON.stringify($scope.Store),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {

            alert(data.data.Message);

            if(data.data.Success==true){
                window.location.href="#!/store_list";
            }
        });

    }
    $scope.uploadImages = function (progressID,file, errFiles) {

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Data;

                    if($scope.Images.indexOf(url)==-1){
                        $scope.Images.push(url);
                    }

                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            UpImageError(errFiles);
        }
    }
    $scope.uploadPictures = function (progressID,file, errFiles) {

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Data;

                    if($scope.Pictures.indexOf(url)==-1){
                        $scope.Pictures.push(url);
                    }


                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            //alert(JSON.stringify(errFiles))
            UpImageError(errFiles);
        }
    }

});
function UpImageError(error){
    var errorTxt ="";
    for(var i=0;i<error.length;i++){

        if(errorTxt==""){
            errorTxt=errorTxt+error[i].$error+":"+error[i].$errorParam;
        }else{
            errorTxt="/"+errorTxt+error[i].$error+":"+error[i].$errorParam;
        }
    }
    if(errorTxt!="" && errorTxt!=undefined){
        alert(errorTxt);
    }

}
main.controller("voucher_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.Voucher =null;
    $scope.TargetAction=null;
    var table;

    $scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify({Name:$scope.Voucher.Name,Amount:$scope.Voucher.Amount,UseDay:$scope.Voucher.UseDay,Introduce:$scope.Voucher.Introduce}),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.Voucher =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }
    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Amount",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"UseDay",render:function (data) {
                        return data+"天";
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "voucher/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            $http.get("voucher/"+row.data().ID,{}).then(function (data) {

                $timeout(function () {
                    $scope.Voucher=data.data.Data;
                    $scope.showModal({title:'修改卡券',url:'voucher/'+row.data().ID,method:'PUT'});
                });

            });




        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("voucher/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );
        
    });
})
main.controller("fullcut_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.FullCut =null;
    $scope.TargetAction=null;
    var table;

    $scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify($scope.FullCut),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.FullCut =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }
    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Amount",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CutAmount",render:function (data) {
                        return $filter("currency")(data/100);
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
           "order":[[1,"asc"]],
            "ajax": {
                "url": "fullcut/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            $http.get("fullcut/"+row.data().ID,{}).then(function (data) {

                $timeout(function () {
                    $scope.FullCut=data.data.Data;
                    $scope.showModal({title:'修改满减',url:'fullcut/save',method:'POST'});
                });

            });
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("fullcut/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})
main.controller("collage_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {



    $scope.Item =null;
    $scope.TargetAction=null;
    let table;


    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Num"},
                {data:"Discount",render:function (data, type, row) {
                        return row.Discount+"%";

                    }},
                {data:"TotalNum"},
                //{data:"GoodsID"},
                {data:"Hash",visible:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改/查看</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
           "order":[[1,"asc"]],
            "ajax": {
                "url": "collage/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    return JSON.stringify(d);
                }
            }
        });



        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());
            window.location.href="#!/add_collage?Hash="+row.data().Hash;
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("collage/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})
main.controller("add_collage_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {



    $scope.Item =null;
    $scope.GoodsList =[];


    var goods_list_table;


    if($routeParams.Hash!=undefined){

        $scope.TargetAction={title:'修改拼团',url:'collage/change',method:'POST'};

        $http.get("collage/"+$routeParams.Hash,{}).then(function (data) {

            var Item = data.data.Data;
            Item.StartTime=new Date(Item.StartTime);
            $scope.Item=Item;
            //$scope.showModal({title:'修改优惠券',url:'timesell/save',method:'POST'});
            //timesell/goods/:TimeSellID/list
            //$scope.listTimeSellGoods();

        });

    }else{
        $scope.TargetAction={title:'添加拼团',url:'collage/save',method:'POST'};
    }
    /*$scope.listTimeSellGoods = function(){
        $http.get("collage/goods/"+$routeParams.Hash+"/list",{}).then(function (data) {
            $scope.GoodsList = data.data.Data;
        });
    }*/

    /*$scope.deleteTimeSellGoods = function(m){

        //timesell/goods/:GoodsID

        if(confirm("是否要取消这个产品的拼团？")){
            $http.delete("collage/goods/"+m.ID,{}).then(function (data) {
                alert(data.data.Message);
                $scope.listTimeSellGoods();
            });
        }


    }*/


    //#!/add_timesell


    $scope.add_score_goods = function(){

        /*if($scope.GoodsList.length==0){
            alert("请先添加产品");
            return
        }*/

        var form ={};
        form.Collage=JSON.stringify($scope.Item);
        /*var GoodsListIDs =[];
        for(var i=0;i<$scope.GoodsList.length;i++){
            GoodsListIDs.push($scope.GoodsList[i].ID);
        }
        form.GoodsListIDs=JSON.stringify(GoodsListIDs);*/
        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {


            alert(data.data.Message);
            alert("前往限时抢购商品管理页面，管理商品");
            window.location.href="#!/collage_manager?Hash="+data.data.Data.Hash;
            $scope.Item =null;
            $scope.TargetAction=null;

            //window.history.back();

        });

    }
   /* $scope.showGoodsList=function(){
        $("#goods_list").modal("show");

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }*/
    $timeout(function () {

        /*goods_list_table = $('#goods_list_table').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:"TimeSellID",className:"opera",orderable:false,render:function (data, type, row) {
                        console.log("--------",row);

                        var have = false;
                        for(var i=0;i<$scope.GoodsList.length;i++){
                            var mitem = $scope.GoodsList[i];
                            if(row.ID==mitem.ID){
                                have=true;
                                break;
                            }
                        }

                        if(have){
                            return '<button class="ui mini button">已经选择</button>';
                        }else{
                            return '<button class="ui select blue mini button">添加</button>';
                        }
                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                //"url": "goods?action=list_goods",
                "url": "goods?action=collage_goods",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#goods_list_table').on('click','td.opera .select', function () {


            var tr = $(this).closest('tr');
            var row = goods_list_table.row( tr );

            console.log(row.data());

            var itme = row.data();
            var have = false;
            for(var i=0;i<$scope.GoodsList.length;i++){
                var mitem = $scope.GoodsList[i];
                if(itme.ID==mitem.ID){
                    have=true;
                    break;
                }
            }

            if(have==false){
                $scope.$apply(function () {
                    $scope.GoodsList.push(itme);
                });

            }

            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList
            goods_list_table.draw(false);
        });*/

    });
});
main.controller("timesell_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.H =[];
    for(var i=0;i<24;i++){
        $scope.H.push({k:i,v:i});
    }
    $scope.M =[];
    for(var i=0;i<60;i++){
        $scope.M.push({k:i,v:i});
    }

    $scope.Item =null;
    $scope.TargetAction=null;
    var table;
    var goods_table;

    /*$scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify($scope.Item),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.Item =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });

    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }*/


    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"BuyNum"},
                {data:"DayNum"},
                {data:"Hash",visible:false},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改/查看</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
           "order":[[1,"asc"]],
            "ajax": {
                "url": "timesell/datatables/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function (d){
                    return JSON.stringify(d);
                }
            }
        });



        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());
            window.location.href="#!/add_timesell?Hash="+row.data().Hash;
        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("timesell/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})
main.controller("collage_manager_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    var goods_list_table;
    var TimeSellGoodsList;



    if($routeParams.Hash==undefined){
        alert("参数不足，无法操作");
        return
    }

    /* $scope.listTimeSellGoods = function(){
         $http.get("timesell/goods/"+$routeParams.Hash+"/list",{}).then(function (data) {
             $scope.GoodsList = data.data.Data;
         });
     }
     $scope.listTimeSellGoods();*/
    $scope.showGoodsList=function(){
        $("#goods_list").modal("show");
        goods_list_table.ajax.reload();
    }

    $timeout(function () {

        goods_list_table = $('#goods_list_table').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                        return '<button class="ui select blue mini button">添加</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                //"url": "goods?action=list_goods",
                "url": "goods?action=activity_goods&Hash="+$routeParams.Hash,
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#goods_list_table').on('click','td.opera .select', function () {


            var tr = $(this).closest('tr');
            var row = goods_list_table.row( tr );

            console.log(row.data());

            var itme = row.data();


            var form ={};
            form.GoodsID=itme.ID;
            form.CollageHash=$routeParams.Hash;
            /* $http.post("timesell/goods/add",{}).then(function (data) {
                alert(data.data.Message);
                //$scope.listTimeSellGoods();
            });*/

            $http.post("collage/goods/add",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {


                if(data.data.Success==true){


                }else{
                    alert(data.data.Message);
                }
                goods_list_table.draw(false);
                TimeSellGoodsList.draw(false);
                //table.ajax.reload();

            });

            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList

        });




        //goods_list_table
        TimeSellGoodsList = $('#TimeSellGoodsList').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button class="ui delete blue mini button">删除这个商品</button>';
                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "collage/goods/"+$routeParams.Hash+"/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#TimeSellGoodsList').on('click','td.opera .delete', function () {


            var tr = $(this).closest('tr');
            var row = TimeSellGoodsList.row( tr );

            console.log(row.data());

            var itme = row.data();



            if(confirm("是否要取消这个产品的限时抢购？")){
                $http.delete("collage/goods/"+itme.ID,{}).then(function (data) {
                    alert(data.data.Message);
                    //$scope.listTimeSellGoods();
                    TimeSellGoodsList.draw(false);
                });
            }


            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList
            //TimeSellGoodsList.draw(false);
        });

    });
});
main.controller("timesell_manager_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    var goods_list_table;
    var TimeSellGoodsList;



    if($routeParams.Hash==undefined){
        alert("参数不足，无法操作");
        return
    }

   /* $scope.listTimeSellGoods = function(){
        $http.get("timesell/goods/"+$routeParams.Hash+"/list",{}).then(function (data) {
            $scope.GoodsList = data.data.Data;
        });
    }
    $scope.listTimeSellGoods();*/
    $scope.showGoodsList=function(){
        $("#goods_list").modal("show");
        goods_list_table.ajax.reload();
    }

    $timeout(function () {

        goods_list_table = $('#goods_list_table').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {

                    return '<button class="ui select blue mini button">添加</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                //"url": "goods?action=list_goods",
                "url": "goods?action=activity_goods&Hash="+$routeParams.Hash,
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#goods_list_table').on('click','td.opera .select', function () {


            var tr = $(this).closest('tr');
            var row = goods_list_table.row( tr );

            console.log(row.data());

            var itme = row.data();


            var form ={};
            form.GoodsID=itme.ID;
            form.TimeSellHash=$routeParams.Hash;
            /* $http.post("timesell/goods/add",{}).then(function (data) {
                alert(data.data.Message);
                //$scope.listTimeSellGoods();
            });*/

            $http.post("timesell/goods/add",$.param(form), {
                transformRequest: angular.identity,
                headers: {"Content-Type": "application/x-www-form-urlencoded"}
            }).then(function (data, status, headers, config) {


                if(data.data.Success==true){


                }else{
                    alert(data.data.Message);
                }
                goods_list_table.draw(false);
                TimeSellGoodsList.draw(false);
                //table.ajax.reload();

            });

            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList

        });




        //goods_list_table
        TimeSellGoodsList = $('#TimeSellGoodsList').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Title"},
                {data:"Stock"},
                {data:"Price",render: function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:"CreatedAt",render: function (data, type, row) {

                        return $filter("date")(data,"medium");
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button class="ui delete blue mini button">删除这个商品</button>';
                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "timesell/goods/"+$routeParams.Hash+"/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });

        $('#TimeSellGoodsList').on('click','td.opera .delete', function () {


            var tr = $(this).closest('tr');
            var row = TimeSellGoodsList.row( tr );

            console.log(row.data());

            var itme = row.data();



            if(confirm("是否要取消这个产品的限时抢购？")){
                $http.delete("timesell/goods/"+itme.ID,{}).then(function (data) {
                    alert(data.data.Message);
                    //$scope.listTimeSellGoods();
                    TimeSellGoodsList.draw(false);
                });
            }


            //$("#goods_list").modal("hide");
            //goods_list_table.ajax.reload();
            //$scope.GoodsList
            //TimeSellGoodsList.draw(false);
        });

    });

});
main.controller("add_timesell_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.H =[];
    for(var i=0;i<24;i++){
        $scope.H.push({k:i,v:i});
    }
    $scope.M =[];
    for(var i=0;i<60;i++){
        $scope.M.push({k:i,v:i});
    }

    $scope.Item ={};
    $scope.GoodsList =[];


    var goods_list_table;



    if($routeParams.Hash!=undefined){

        $scope.TargetAction={title:'修改限时抢购',url:'timesell/change',method:'POST'};

        $http.get("timesell/"+$routeParams.Hash,{}).then(function (data) {

            var Item = data.data.Data;
            Item.StartTime=new Date(Item.StartTime);
            $scope.Item=Item;
            //$scope.showModal({title:'修改优惠券',url:'timesell/save',method:'POST'});
            //timesell/goods/:TimeSellID/list
            //$scope.listTimeSellGoods();

        });

    }else{
        $scope.TargetAction={title:'添加限时抢购',url:'timesell/save',method:'POST'};
    }





    //#!/add_timesell


    $scope.add_score_goods = function(){

        /*if($scope.GoodsList.length==0){
            alert("请先添加产品");
            return
        }*/

        var form ={};
        form.TimeSell=JSON.stringify($scope.Item);
        /* var GoodsListIDs =[];
        for(var i=0;i<$scope.GoodsList.length;i++){
            GoodsListIDs.push($scope.GoodsList[i].ID);
        }*/
        //form.GoodsListIDs=JSON.stringify(GoodsListIDs);
        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {


            alert(data.data.Message);
            ///window.history.back();
            alert("前往限时抢购商品管理页面，管理商品");
            window.location.href="#!/timesell_manager?Hash="+data.data.Data.Hash;
            $scope.Item =null;
            $scope.TargetAction=null;

        });

    }

    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }

});

main.controller("score_goods_list_controller",function ($http,$filter,$scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.ScoreGoods =null;
    $scope.TargetAction=null;
    var table;



    $scope.uploadImages = function (progressID,file, errFiles) {

        if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Data;

                    $scope.ScoreGoods.Image=url;

                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            UpImageError(errFiles);
        }
    }

    $scope.add_score_goods = function(){

        $http({
            method:$scope.TargetAction.method,
            url:$scope.TargetAction.url,
            data:JSON.stringify($scope.ScoreGoods),
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
            $scope.ScoreGoods =null;
            $scope.TargetAction=null;
            alert(data.data.Message);
            $("#add_score_goods").modal("hide");
            table.ajax.reload();
        });
    }
    $scope.showModal = function (ta) {
        $scope.TargetAction = ta;
        $("#add_score_goods").modal("show");
    }
    $timeout(function () {

       table = $('#table_local').DataTable({
            "columns": [
                {data:"ID"},
                {data:"Name"},
                {data:"Score"},
                {data:"Price",render:function (data, type, row) {
                        return $filter("currency")(data/100);
                    }},
                {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                        return '<button  class="ui edit blue mini button">修改</button>'+
                            '<button class="ui delete red mini button">删除</button>';

                    }}
            ],
            "createdRow": function ( row, data, index ) {
                //console.log(row,data,index);
            },
            columnDefs:[

            ],
            "initComplete":function (d) {

            },
            paging: true,
            //"dom": '<"toolbar">frtip',
            "pagingType": "full_numbers",
            searching: true,
            "processing": true,
            "serverSide": true,
            "ajax": {
                "url": "score_goods/list",
                "type": "POST",
                "contentType": "application/json",
                "data": function ( d ) {
                    return JSON.stringify(d);
                }
            }
        });


        $('#table_local').on('click','td.opera .edit', function () {
            var tr = $(this).closest('tr');
            var row = table.row( tr );
            console.log(row.data());


            $http.get("score_goods/"+row.data().ID,{}).then(function (data) {

                $timeout(function () {
                    $scope.ScoreGoods=data.data.Data;
                    $scope.showModal({title:'修改积分产品',url:'score_goods/'+row.data().ID,method:'PUT'});
                });

            });




        });

        $('#table_local').on('click','td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

                if (confirm("确定删除？")) {

                    $http.delete("score_goods/"+row.data().ID,{}).then(function (data) {
                        alert(data.data.Message);
                        table.ajax.reload();
                    });

                }
            }
        );

    });
})



main.controller("main_controller",function ($http, $scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {
    
});


function LeftMoveArr(arr,index){
    var item =arr[index];
    if(index-1<0){
        return
    }
    var newIndex =index-1;

    var oldItem = arr[newIndex];

    arr[newIndex]=item;
    arr[index]=oldItem;

}
function RightMoveArr(arr,index){

    var item =arr[index];
    if(index+1>arr.length-1){
        return
    }
    var newIndex =index+1;

    var oldItem = arr[newIndex];

    arr[newIndex]=item;
    arr[index]=oldItem;

}
main.controller("add_goods_controller",function ($http, $scope, $rootScope, $routeParams,$document,$timeout,$interval,Upload) {

    $scope.LeftMoveArr=LeftMoveArr;
    $scope.RightMoveArr=RightMoveArr;

    $scope.Images=[];
    $scope.Videos=[];
    $scope.Pictures=[];
    $scope.Params=[];

    $scope.Goods ={ID:$routeParams.ID};
    $scope.param = {Name:"",Value:""};
    $scope.GoodsTypeList = [];
    $scope.GoodsTypeID=undefined;

    $scope.Specification={};
    $scope.Specifications=[];

    $scope.PostAction={Action:"POST",Url:"goods?action=add_goods"};


    $scope.addSpecifications = function(){

        if($scope.Specification.Num==undefined){
            alert("请填写：Num");
            return
        }
        if($scope.Specification.Weight==undefined){
            alert("请填写：Weight");
            return
        }
        if($scope.Specification.Label==undefined){
            alert("请填写：Label");
            return
        }
        if($scope.Specification.Stock==undefined){
            alert("请填写：Stock");
            return
        }
        if($scope.Specification.CostPrice==undefined){
            alert("请填写：CostPrice");
            return
        }
        if($scope.Specification.MarketPrice==undefined){
            alert("请填写：MarketPrice");
            return
        }
        if($scope.Specification.Brokerage==undefined){
            alert("请填写：Brokerage");
            return
        }
        var copy = angular.copy($scope.Specification);
        copy.Delete = 0;

        $scope.Specifications.push(copy);
        $scope.Specification={};

        $scope.ayStock();
    }

    $scope.ayStock = function(){
        var stock = 0;
        var Price = 9999999999999999;
        for(var i=0;i<$scope.Specifications.length;i++){
            var item = $scope.Specifications[i];
            stock=stock+parseInt(item.Stock);
            Price = Math.min(Price,item.MarketPrice);
        }
        $scope.Goods.Stock=stock;
        $scope.Goods.Price=Price;
    }
    $scope.deleteSpecification = function(index){
        var item = $scope.Specifications[index];
        $scope.Specifications.splice(index,1);
        $scope.ayStock();
        if(item.ID!=undefined&& item.ID>0){
            $http.get("goods?action=delete_specification&ID="+item.ID).then(function (data, status, headers, config) {
                alert(data.data.Message);
            });
        }


    }

    $scope.changeStock = function(){
        $scope.ayStock();
    }


    $scope.expressTemplateInfo=null;
    $scope.selectExpressTemplate = function(){
        $scope.Units=[];
        for(var i=0;i<$scope.ExpressTemplateList.length;i++){
            var item = $scope.ExpressTemplateList[i];
            if(item.ID==$scope.Goods.ExpressTemplateID){
                $scope.expressTemplateInfo="当前快递为："+item.Name+"，"+(item.Drawee=='BUYERS'?'买家承担运费':'商家包邮')+"，计费方式："+(item.Type=='ITEM'?'件':'Kg');
                break
            }
        }
    }

    $scope.ExpressTemplateList =[];
    $http.get("express_template/list").then(function (data, status, headers, config) {
        $scope.ExpressTemplateList = data.data.Data;

        $scope.selectExpressTemplate();
    });

    $http.get("goods?action=list_goods_type_all",{}, {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        $scope.GoodsTypeList = data.data.Data;

            if($scope.Goods.ID!=undefined){

                $scope.PostAction={Action:"POST",Url:"goods?action=change_goods"};

                $http.get("goods?action=get_goods",{params:{ID:$scope.Goods.ID}}).then(function (data) {

                    if(data.data.Success==true){

                        var Goods = data.data.Data.Goods;
                        Goods.Price= Goods.Price/100;
                        $scope.Goods =Goods;

                        var Specifications=data.data.Data.Specifications;
                        for(var i=0;i<Specifications.length;i++){
                            Specifications[i].Weight=Specifications[i].Weight/1000;
                            Specifications[i].CostPrice=Specifications[i].CostPrice/100;
                            Specifications[i].MarketPrice=Specifications[i].MarketPrice/100;
                            Specifications[i].Brokerage=Specifications[i].Brokerage/100;
                        }
                        $scope.Specifications=Specifications;

                        $scope.Videos=JSON.parse($scope.Goods.Videos);
                        $scope.Images=JSON.parse($scope.Goods.Images);
                        $scope.Pictures=JSON.parse($scope.Goods.Pictures);
                        $scope.Params=JSON.parse($scope.Goods.Params);

                        $scope.selectExpressTemplate();

                        $http.get("goods?action=get_goods_type_child",{params:{ID:$scope.Goods.GoodsTypeChildID}}, {
                            transformRequest: angular.identity,
                            headers: {"Content-Type": "application/json"}
                        }).then(function (data, status, headers, config) {
                            var GoodsTypeChild = data.data.Data;


                            $http.get("goods?action=list_goods_type_child_id",{params:{ID:GoodsTypeChild.GoodsTypeID}}, {
                                transformRequest: angular.identity,
                                headers: {"Content-Type": "application/json"}
                            }).then(function (data, status, headers, config) {
                                $scope.GoodsTypeChildList = data.data.Data;

                                $timeout(function () {
                                    $scope.GoodsTypeChildID =GoodsTypeChild.ID;
                                    $scope.GoodsTypeID=GoodsTypeChild.GoodsTypeID;
                                });

                            });




                        });




                    }else{
                        alert(data.data.Message);
                    }

                })
            }





    });


    $scope.changeGoodsType = function(){
        $scope.GoodsTypeChildList=[];

        if($scope.GoodsTypeID!=undefined){
            $http.get("goods?action=list_goods_type_child_id",{params:{ID:$scope.GoodsTypeID}}).then(function (data) {

                $scope.GoodsTypeChildList = data.data.Data;

            })
        }
    }


    $scope.deleteArr = function(arr,index){

        if(confirm("确认删除这项内容？")){
            arr.splice(index,1);
        }
    }
    $scope.showParamsModal = function(){
        $('#params').modal({
            onApprove : function() {
                window.alert('Approved!');
            }
        }).modal('show');
    }

    $scope.addParams = function(){

        $timeout(function () {
            $scope.Params.push(angular.copy($scope.param));
            $('#params').modal("hide");
            $scope.param = {Name:"",Value:""};
        });
    }




    $scope.saveGoods = function(){

        $scope.Goods.Videos = JSON.stringify($scope.Videos);
        $scope.Goods.Images = JSON.stringify($scope.Images);
        $scope.Goods.Pictures = JSON.stringify($scope.Pictures);
        $scope.Goods.Params = JSON.stringify($scope.Params);
        //$scope.Goods.Specifications = $scope.Specifications;
        $scope.Goods.GoodsTypeID = parseInt($scope.GoodsTypeID);
        $scope.Goods.GoodsTypeChildID = parseInt($scope.GoodsTypeChildID);
        $scope.Goods.Price =$scope.Goods.Price*100;//parseInt($scope.Goods.Price*100);
        /*var form = new FormData();
        form.append("goods",JSON.stringify($scope.Goods));
        form.append("specifications",JSON.stringify($scope.Specifications));*/



        if( $scope.Specifications.length<=0){
            alert("请添加规格");
            return;
        }





        var Specifications=$scope.Specifications;

        for(var i=0;i<Specifications.length;i++){
            Specifications[i].Weight=parseInt(Specifications[i].Weight*1000);
            Specifications[i].CostPrice=parseInt(Specifications[i].CostPrice*100);
            Specifications[i].MarketPrice=parseInt(Specifications[i].MarketPrice*100);
            Specifications[i].Brokerage=parseInt(Specifications[i].Brokerage*100);
        }
        $scope.Specifications=Specifications;


        var form ={};
        form.goods=JSON.stringify($scope.Goods);
        form.specifications=JSON.stringify($scope.Specifications);
        //$scope.PostAction={Action:"POST",Url:"goods?action=change_goods"};
        $http({
            method:$scope.PostAction.Action,
            url:$scope.PostAction.Url,
            data:$.param(form),
            transformRequest: angular.identity,
            headers: {'Content-Type':'application/x-www-form-urlencoded'}
        }).then(function (data, status, headers, config) {
            alert(data.data.Message);
            if(data.data.Success==true){
                window.location.href="#!/goods_list";
            }
        });

        /*$http.post("goods?action="+action,$.param(form), {
            transformRequest: angular.identity,
            //headers: {"Content-Type": "application/x-www-form-urlencoded"}
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
        }).then(function (data, status, headers, config) {

           alert(data.data.Message);
            if(data.data.Success==true){
                window.location.href="#!/goods_list";
            }
        });*/

    }
    $scope.uploadVideos = function (progressID,files, errFiles) {

        if (files && files.length) {

            //progress-bar-videos
            var progressObj ={};
            $(progressID).text("0/"+(files.length*100));

            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Data;

                    if($scope.Videos.indexOf(url)==-1){
                        $scope.Videos.push(url);
                    }

                },function (resp) {
                    console.log('Error status: ' + resp.status);
                },function (evt) {
                    var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
                    console.log('progress: ' + progressPercentage + '% ' + evt.config.data.file.name);


                    progressObj[evt.config.data.file.name]=progressPercentage;


                    var showTexts =[];
                    for(var key in progressObj){
                        showTexts.push(key+":"+progressObj[key]+"%");
                    }


                    $(progressID).text(showTexts.join(","));
                });
            }
        }else{
            UpImageError(errFiles);
        }
    }
    $scope.uploadImages = function (progressID,files, errFiles) {

        if (files && files.length) {
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Data;

                    if($scope.Images.indexOf(url)==-1){
                        $scope.Images.push(url);
                    }
                    
                },function (response) {
                    
                },function (response) {
                    
                });
            }
        }else{
            UpImageError(errFiles);
        }
    }
    $scope.uploadPictures = function (progressID,files, errFiles) {

        if (files && files.length) {
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up',data:{file: files[i]}}).then(function (response) {
                    var url =response.data.Data;
                    if($scope.Pictures.indexOf(url)==-1){
                        $scope.Pictures.push(url);
                    }

                },function (response) {

                },function (response) {

                });
            }
        }else{
            UpImageError(errFiles);
        }

        /*if (file) {
            var thumbnail =Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                $timeout(function () {
                    var url =response.data.Data;

                    if($scope.Pictures.indexOf(url)==-1){
                        $scope.Pictures.push(url);
                    }


                });
            }, function (response) {
                if (response.status > 0){
                    $scope.errorMsg = response.status + ': ' + response.data;
                }
            }, function (evt) {
                // Math.min is to fix IE which reports 200% sometimes
                //var progress = Math.min(100, parseInt(100.0 * evt.loaded / evt.total));
                //$("."+progressID).text(progress+"%");
                //$("."+progressID).css("width",progress+"%");
            });
        }else{
            UpImageError(errFiles);
        }*/
    }

});
main.controller("goods_list_controller",function ($http, $scope, $filter,$rootScope, $routeParams,$document,$timeout,$interval,Upload) {


    $http.get("goods?action=list_goods_type_child",{}, {
        transformRequest: angular.identity,
        headers: {"Content-Type": "application/json"}
    }).then(function (data, status, headers, config) {
        var _list = data.data.Data;
        var GoodsTypeObj = {};
        for(var i=0;i<_list.length;i++){
            GoodsTypeObj[_list[i].ID]=_list[i].Name;
        }



        var table;
        $timeout(function () {
            table = $('#table_local').DataTable({
                "columns": [
                    {data:"ID"},
                    {data:"Title"},
                    {data:"Stock"},
                   /* {data:"CostPrice",render: function (data, type, row) {

                            return $filter("currency")(data/100);
                        }},
                    {data:"MarketPrice",render: function (data, type, row) {

                            return $filter("currency")(data/100);
                        }},*/
                    {data:"CreatedAt",render: function (data, type, row) {

                            return $filter("date")(data,"medium");
                        }},
                    {data:"GoodsTypeChildID",render: function (data, type, row) {
                            if(GoodsTypeObj[data]){
                                return GoodsTypeObj[data];
                            }else {
                                return "系列不存在"
                            }

                        }},
                    {data:null,className:"opera",orderable:false,render:function (data, type, row) {
                            return '<a href="#!/add_goods?ID='+data.ID+'" class="ui edit blue mini button">编辑</a>'+
                                '  <button class="ui delete red mini button">删除</button>';

                        }}
                ],
                "createdRow": function ( row, data, index ) {
                    //console.log(row,data,index);
                },
                columnDefs:[

                ],
                "initComplete":function (d) {

                },
                paging: true,
                //"dom": '<"toolbar">frtip',
                "pagingType": "full_numbers",
                searching: false,
                "processing": true,
                "serverSide": true,
                "ajax": {
                    "url": "goods?action=list_goods",
                    "type": "POST",
                    "contentType": "application/json",
                    "data": function ( d ) {
                        return JSON.stringify(d);
                    }
                }
            });

            /*$('#table_local').on('click','td.opera .edit', function () {


                var tr = $(this).closest('tr');
                var row = table.row( tr );
                console.log(row.data());

                $timeout(function () {
                    $scope.GoodsType={Name:row.data().Name,ID:row.data().ID};
                    $scope.showGoodsTypeModal(1);
                });

            });*/
            $('#table_local').on('click','td.opera .delete', function () {


                var tr = $(this).closest('tr');
                var row = table.row( tr );

                console.log(row.data());

                if(confirm("确定删除？")){
                    $http.get("goods?action=del_goods",{params:{ID:row.data().ID}}).then(function (data) {

                        alert(data.data.Message);

                        table.ajax.reload();

                    })
                }

                /*$timeout(function () {
                    var data = row.data();
                    data.PassWord="";
                    $scope.onShowBox(data,1);
                });*/



            });
        });
    })



});
