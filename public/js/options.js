/* global GLOBAL */
$(function () {
    var valToId = [];
    valToId.push({
        key: "radius",
        id: "size-radius",
        isNumber: true
    });
    valToId.push({
        key: "fontSize",
        id: "size-font",
        isNumber: true
    });
    valToId.push({
        key: "width",
        id: "size-width",
        isNumber: true
    });
    valToId.push({
        key: "height",
        id: "size-height",
        isNumber: true
    });
    valToId.push({
        key: "fontColor",
        id: "font-color"
    });
    valToId.push({
        key: "backColor",
        id: "background-color"
    });
    valToId.push({
        key: "nodeColor",
        id: "node-color"
    });
    valToId.push({
        key: "nodeBorderColor",
        id: "node-border-color"
    });
    valToId.push({
        key: "arrowColor",
        id: "arrow-color"
    });
    valToId.push({
        key: "arrowColor",
        id: "arrow-color"
    });

    var options = {};

    options.radius = 15;
    options.fontSize = 17;
    options.fontColor = "#000000";
    options.nodeBorderColor = "#000000";
    options.nodeColor = "#ffffff";
    options.arrowColor = "#000000";
    options.backColor = "#ffffff";
    options.width = 250;
    options.height = 400;

    var $canvasContainer = $('.canvas-container')

    options.width = $canvasContainer.width();
    options.height = $canvasContainer.height();

    function reset() {
        for (var i = 0; i < valToId.length; ++i) {
            $('#' + valToId[i].id).val(options[valToId[i].key])
        }
    }

    reset();

    function init() {
        function initChange(id, key, type) {
            $('#' + id).change(function () {
                options[key] = $(this).val();

                if (type) {
                    options[key] = parseInt(options[key])
                }

                GLOBAL.redraw();
            });
        }

        for (var i = 0; i < valToId.length; ++i) {
            initChange(valToId[i].id, valToId[i].key, valToId[i].isNumber)
        }

        GLOBAL.options = options;
    }

    init();

});