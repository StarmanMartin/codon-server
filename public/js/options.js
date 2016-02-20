/* global redraw */
var optoins = {};

optoins.radius = 15;
optoins.fontSize = 17;
optoins.fontColor = "#000000";
optoins.nodeBorderColor = "#000000";
optoins.nodeColor = "#ffffff";
optoins.arrowColor = "#000000";
optoins.backColor = "#ffffff";
optoins.width = 250;
optoins.height = 400;

$(function(){    
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
    var $canvasContainer = $('.canvas-container')
                    
    optoins.width = $canvasContainer.width();
    optoins.height = $canvasContainer.height();
    
    function reset(){
        for(var i = 0; i < valToId.length; ++i) {
            $('#' + valToId[i].id).val(optoins[valToId[i].key]) 
        }
    }
    
    reset();
    
    function init(){
        function initChange(id, key, type) {
            $('#' + id).change(function(){
                optoins[key] = $(this).val();
                
                if (type) {
                    optoins[key] = parseInt(optoins[key])
                }
                
                redraw();
            });
        }
        
        for(var i = 0; i < valToId.length; ++i) {
             initChange(valToId[i].id, valToId[i].key, valToId[i].isNumber)
        }
    }
    
    init();
   
});