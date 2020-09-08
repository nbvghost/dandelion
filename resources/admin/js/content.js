main.config(function ($routeProvider, $locationProvider, $provide, $httpProvider, $httpParamSerializerJQLikeProvider, $interpolateProvider) {

    $routeProvider.when("/content_item_list", {
        templateUrl: "content_templates/content_list_template",
        controller: "content_list_controller"
    });
    $routeProvider.when("/add_contents", {
        templateUrl: "content_templates/add_articles_template",
        controller: "content_add_articles_controller"
    });
    $routeProvider.when("/add_gallery", {
        templateUrl: "content_templates/add_gallery_template",
        controller: "content_add_gallery_controller"
    });
    $routeProvider.when("/contents", {
        templateUrl: "content_templates/articles_template",
        controller: "content_articles_controller"
    });
    $routeProvider.when("/gallery", {
        templateUrl: "content_templates/articles_template",
        controller: "content_articles_controller"
    });
    $routeProvider.when("/content", {
        templateUrl: "content_templates/add_article_template",
        controller: "content_add_article_controller"
    });


});


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
            window.location.href = "#!/add_contents?ContentItemID=" + row.data().ContentItemID + "&ID=" + row.data().ID;
        });

        $('#table_local').on('click', 'td.opera .delete', function () {
                var tr = $(this).closest('tr');
                var row = table.row(tr);
                console.log(row.data());

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

main.controller("content_add_articles_controller", function ($http, $scope, $routeParams, $rootScope, $timeout, $location, Upload) {
    $scope.ContentSubTypes = {};

    $scope.Article = {ContentItemID: parseInt($routeParams.ContentItemID)};


    $scope.saveArticle = function () {

        //$scope.ContentSubTypeID;
        //$scope.ContentSubTypeChildID;
        //console.log(quill.container.firstChild.innerHTML)
        $scope.Article.ContentItemID = parseInt($routeParams.ContentItemID);
        $scope.Article.Content = quill.container.firstChild.innerHTML;


        $scope.Article.ContentSubTypeID = parseInt($scope.Article.ContentSubTypeID)

        if (parseInt($scope.Article.ContentSubTypeChildID) > 0) {
            $scope.Article.ContentSubTypeID = parseInt($scope.Article.ContentSubTypeChildID)
        }

        if (!$scope.Article.ContentSubTypeID) {
            //$scope.Article.ContentSubTypeID=$scope.ContentSubTypeID;
            //$scope.Article.ContentSubTypeChildID=$scope.ContentSubTypeChildID;
            //}else{
            alert("请选择分类");
            return
        }
        $http.post("content/save", JSON.stringify($scope.Article), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {
            console.log(data);
            alert(data.data.Message);
            if (data.data.Success) {
                window.history.back();
            }
        });
    }


    $http.get("content/sub_type/list/all/" + $routeParams.ContentItemID).then(function (value) {

        $scope.ContentSubTypes = value.data.Data || [];
    });


    $scope.changeContentSubTypes = function () {
        //$scope.ContentSubTypeChildID=undefined;
        //$scope.Article.ContentSubTypeID=$scope.MContentSubTypeID;
        //$scope.ContentSubTypeChilds=[];
        //console.log($scope.ContentSubTypeID);
        //$scope.listContentSubTypeChilds($routeParams.ContentItemID,$scope.Article.ContentSubTypeID);
    }
    $scope.changeContentSubTypeChilds = function () {
        //alert($scope.MContentSubTypeChildID);
        /* if(parseInt($scope.MContentSubTypeChildID)>0){
             $scope.Article.ContentSubTypeID=parseInt($scope.MContentSubTypeChildID)
         }else{
             $scope.Article.ContentSubTypeID=parseInt($scope.MContentSubTypeID);
         }*/

        //$scope.Article.;
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
    let quill;
    $timeout(function () {

        const Inline = Quill.import('blots/inline');
        const Block = Quill.import('blots/block');
        const BlockEmbed = Quill.import('blots/block/embed');

        class BoldBlot extends Inline {
        }

        BoldBlot.blotName = 'bold';
        BoldBlot.tagName = 'strong';

        class ItalicBlot extends Inline {
        }

        ItalicBlot.blotName = 'italic';
        ItalicBlot.tagName = 'em';

        class LinkBlot extends Inline {
            static create(url) {
                const node = super.create();
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

        class BlockquoteBlot extends Block {
        }

        BlockquoteBlot.blotName = 'blockquote';
        BlockquoteBlot.tagName = 'blockquote';

        class HeaderBlot extends Block {
            static formats(node) {
                return HeaderBlot.tagName.indexOf(node.tagName) + 1;
            }
        }

        HeaderBlot.blotName = 'header';
        HeaderBlot.tagName = ['H1', 'H2'];

        class DividerBlot extends BlockEmbed {
        }

        DividerBlot.blotName = 'divider';
        DividerBlot.tagName = 'hr';

        class ImageBlot extends BlockEmbed {
            static create(value) {
                const node = super.create();
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
                const node = super.create();
                node.setAttribute('src', url);
                node.setAttribute('frameborder', '0');
                node.setAttribute('allowfullscreen', true);
                return node;
            }

            static formats(node) {
                const format = {};
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

        if ($routeParams.ID) {
            $http.get("content/multi/get/" + $routeParams.ID).then(function (responea) {

                $scope.Article = responea.data.Data;
                quill.clipboard.dangerouslyPasteHTML($scope.Article.Content);


                if ($scope.Article.ContentSubTypeID > 0) {
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

                        });


                    });
                }


            });
        }

        quill.getModule("toolbar").addHandler("image", function (e) {

            //var baseUrl ="//"+$location.host()+":"+$location.port();

            $("#SelectImageModal").modal({
                onApprove: function (e) {


                    if ($scope.EditImages.length > 0) {


                        for (let ii = 0; ii < $scope.EditImages.length; ii++) {

                            const range = quill.getSelection(true);
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
                            $scope.EditImages = [];
                        });
                        return true;


                    } else {
                        return false;
                    }
                }, closable: false
            }).modal("show");
        });

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
main.controller("content_add_article_controller", function ($http, $scope, $routeParams, $rootScope, $timeout, $location, Upload) {


    let quill;
    $scope.ContentSubTypes = [];
    $scope.ContentItemID = $routeParams.ContentItemID

    $scope.ContentSubTypeID = $routeParams.ContentSubTypeID
    $scope.ContentSubTypeChildID = $routeParams.ContentSubTypeChildID


    $scope.Article = {ContentItemID: $scope.ContentItemID};


    $scope.saveArticle = function () {

        //$scope.ContentSubTypeID;
        //$scope.ContentSubTypeChildID;
        //console.log(quill.container.firstChild.innerHTML)
        $scope.Article.ContentItemID = parseInt($scope.ContentItemID);

        $scope.Article.Content = quill.container.firstChild.innerHTML;
        $scope.Article.ContentSubTypeID = parseInt($scope.Article.ContentSubTypeChildID);

        $http.post("content/save", JSON.stringify($scope.Article), {
            transformRequest: angular.identity,
            headers: {"Content-Type": "application/json"}
        }).then(function (data, status, headers, config) {
            console.log(data);
            alert(data.data.Message);
            if (data.data.Success) {
                window.history.back();
            }
        });
    }

    $http.get("content/sub_type/list/all/" + $routeParams.ContentItemID).then(function (value) {
        $scope.ContentSubTypes = value.data.Data || [];

        $scope.Article.ContentSubTypeID=$scope.ContentSubTypeID
        $scope.Article.ContentSubTypeChildID=$scope.ContentSubTypeChildID

    });

    $scope.UploadPictureImage = function (file, errFiles) {
        if (file) {
            var thumbnail = Upload.upload({
                url: '/file/up',
                data: {file: file},
            });
            thumbnail.then(function (response) {
                var url = response.data.Path;
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
            for (var i = 0; i < files.length; i++) {
                Upload.upload({url: '/file/up', data: {file: files[i]}}).then(function (response) {
                    var url = response.data.Path;

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
    $timeout(function () {

        const Inline = Quill.import('blots/inline');
        const Block = Quill.import('blots/block');
        const BlockEmbed = Quill.import('blots/block/embed');

        class BoldBlot extends Inline {
        }

        BoldBlot.blotName = 'bold';
        BoldBlot.tagName = 'strong';

        class ItalicBlot extends Inline {
        }

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

        class BlockquoteBlot extends Block {
        }

        BlockquoteBlot.blotName = 'blockquote';
        BlockquoteBlot.tagName = 'blockquote';

        class HeaderBlot extends Block {
            static formats(node) {
                return HeaderBlot.tagName.indexOf(node.tagName) + 1;
            }
        }

        HeaderBlot.blotName = 'header';
        HeaderBlot.tagName = ['H1', 'H2'];

        class DividerBlot extends BlockEmbed {
        }

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

        if ($scope.ContentItemID && $scope.ContentSubTypeID && $scope.ContentSubTypeChildID) {

            $http.get("content/single/get/" + $scope.ContentItemID+"/"+$scope.ContentSubTypeChildID).then(function (responea) {

                if(responea.data.Data.ID>0){
                    $scope.Article = responea.data.Data;

                    if( $scope.Article.ContentSubTypeID!==parseInt($scope.ContentSubTypeChildID)){

                        alert("原内容类型与所选类型不匹配")
                        return
                    }
                    //$scope.ContentSubTypeID = $scope.Article.ContentSubTypeID
                    //$scope.ContentSubTypeChildID = $scope.Article.ContentSubTypeChildID
                    $scope.Article.ContentSubTypeID=$scope.ContentSubTypeID
                    $scope.Article.ContentSubTypeChildID=$scope.ContentSubTypeChildID

                    quill.clipboard.dangerouslyPasteHTML($scope.Article.Content);
                }

            });
        }

        quill.getModule("toolbar").addHandler("image", function (e) {

            //var baseUrl ="//"+$location.host()+":"+$location.port();

            $("#SelectImageModal").modal({
                onApprove: function (e) {


                    if ($scope.EditImages.length > 0) {


                        for (var ii = 0; ii < $scope.EditImages.length; ii++) {

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
                            $scope.EditImages = [];
                        });
                        return true;


                    } else {
                        return false;
                    }
                }, closable: false
            }).modal("show");
        });

    });


});

main.controller('content_list_controller', function ($http, $scope, $timeout, $routeParams, $document, $interval) {

    $scope.MenuTypes = [];
    $scope.EditMenus=null;
    $scope.CreateMenus=null;
    $scope.selectClassify = null;
    $scope.classifyChild = {};

    $scope.ContentContentSubType ={}
    $scope.ContentSubTypes ={}

    $scope.templateNameObj = {
        "contents": [
            {Key: "news", Label: "新闻中心", SubMenu: true, Content: true}
        ],
        "content": [
            {Key: "services", Label: "服务", SubMenu: true, Content: true},
            {Key: "about", Label: "关于我们", SubMenu: false, Content: true}
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

            }).catch((error)=>{
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

        $scope.changeMenuSort(current,targetIndex,target);

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


        $scope.changeMenuSort(current,targetIndex,target);

    }
    $scope.changeMenuSort = function (current,targetIndex,target){
        $scope.MenusList[targetIndex] = current;
        $scope.MenusList[index] = target;

        current.Sort = targetIndex;
        target.Sort = index;

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
        $scope.CreateMenus=null;
    }
    $scope.saveEditMenu = async function () {
        ActionTarget = {method: 'PUT', url: 'content/item/' + m.ID, title: '修改菜单'};
        await $scope.saveMenu($scope.EditMenus);
        $scope.EditMenus=null;
    }
    //{method:'PUT',url:''}



    $scope.showContentContentSubTypeDialog = function (m) {


        $http.get("content/sub_type/list/all/" + m.ID).then(function (value) {

            $scope.ContentSubTypes = value.data.Data || [];

            $("#contentContentSubTypeDialog").modal({
                centered: false,closable:false, allowMultiple: false,
                onDeny:function(e){
                    $timeout(function (){
                        $scope.ContentSubTypes={};
                        $scope.ContentContentSubType={};
                    })
                    return true
                },
                onApprove: function (e) {

                    if($scope.ContentContentSubType.ContentSubTypeID==undefined || $scope.ContentContentSubType.ContentSubTypeChildID==undefined){
                        alert("请选择类别")
                        return false
                    }

                    $timeout(function (){
                        let redirect ="#!/"+m.Type+"?ContentItemID="+m.ID+"&ContentSubTypeID="+$scope.ContentContentSubType.ContentSubTypeID+"&ContentSubTypeChildID="+$scope.ContentContentSubType.ContentSubTypeChildID
                        $scope.ContentSubTypes={};
                        $scope.ContentContentSubType={};
                        window.location.href=redirect
                    })

                    return true
                }
            }).modal("show");
        });


    }
    $scope.editMenus = function (m) {
        //$scope.selectClassify = null;
        //$scope.classifyChild = null;
        $scope.EditMenus = m;
        $scope.EditMenus.SubMenu = $scope.getTemplateNameObj($scope.EditMenus.Type,$scope.EditMenus.TemplateName);

        $scope.classify = {ContentItemID: $scope.EditMenus.ID};
        $("#editMenus").modal({
            centered: false,closable:false, allowMultiple: false,
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
        if(confirm("确定要删除？")==false){
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

        if(confirm("确定要删除？")==false){
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
        if(confirm("确定要删除？")==false){
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
