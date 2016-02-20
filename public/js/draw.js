/* global GLOBAL */

$(function () {
    function drawArrow(ctx, nodeA, nodeB) {
        var fromx = nodeA.center.x;
        var fromy = nodeA.center.y;
        var tox = nodeB.center.x;
        var toy = nodeB.center.y;
        var headlen = (GLOBAL.options.radius * 6) / (17);   // length of head in pixels
        var angle = Math.atan2(toy - fromy, tox - fromx);

        fromx += (5 + GLOBAL.options.radius) * Math.cos(angle);
        fromy += (5 + GLOBAL.options.radius) * Math.sin(angle);
        tox -= (5 + GLOBAL.options.radius) * Math.cos(angle); //newPoints[2]
        toy -= (5 + GLOBAL.options.radius) * Math.sin(angle);

        ctx.save();
        ctx.beginPath();
        ctx.strokeStyle = GLOBAL.options.arrowColor;
        ctx.moveTo(fromx, fromy);
        ctx.lineTo(tox, toy);
        ctx.lineTo(tox - headlen * Math.cos(angle - (Math.PI * GLOBAL.options.radius) / (17 * 10)), toy - headlen * Math.sin(angle - (Math.PI * GLOBAL.options.radius) / (6 * 17)));
        ctx.lineTo(tox - headlen * Math.cos(angle + (Math.PI * GLOBAL.options.radius) / (17 * 10)), toy - headlen * Math.sin(angle + (Math.PI * GLOBAL.options.radius) / (6 * 17)));
        ctx.lineTo(tox, toy);
        ctx.stroke();
        ctx.restore();
    }

    GLOBAL.drawGraph = function drawGraph(nodeObj, data) {
        var c = document.getElementById("mycanvas");
        var ctx = c.getContext("2d");
        var width = GLOBAL.options.width;
        var height = GLOBAL.options.height;
        ctx.fillStyle = GLOBAL.options.backColor;
        ctx.fillRect(0, 0, width, height)
        ctx.font = GLOBAL.options.fontSize + "px Arial";
        var nodes = nodeObj.DinucleotideNodes.concat(nodeObj.TetranucleotideNodes)

        for (var i = 0; i < nodes.length; ++i) {
            ctx.fillStyle = GLOBAL.options.nodeColor;
            ctx.strokeStyle = GLOBAL.options.nodeBorderColor
            ctx.beginPath();
            ctx.arc(nodes[i].center.x, nodes[i].center.y, GLOBAL.options.radius, 0, 2 * Math.PI);
            ctx.stroke();
            ctx.arc(nodes[i].center.x, nodes[i].center.y, GLOBAL.options.radius, 0, 2 * Math.PI);
            ctx.fill();

            ctx.fillStyle = GLOBAL.options.fontColor;

            var pos = ctx.measureText(nodes[i].text).width / 2;
            ctx.fillText(nodes[i].text, nodes[i].center.x - pos, ctx.fillText(nodes[i].text, nodes[i].center.x - pos, nodes[i].center.y + GLOBAL.options.fontSize * 0.35));
        }

        var indexTetranucleotide = 0;
        for (var right = 0; right < 7; right += 2) {
            if (data.TetranucleotideNodes[right].length > 0 ||
                data.TetranucleotideNodes[right + 1].length > 0) {
                for (var i = 0; i < data.TetranucleotideNodes[right].length; ++i) {
                    var index = data.TetranucleotideNodes[right][i];
                    drawArrow(ctx, nodeObj.DinucleotideNodes[index],
                        nodeObj.TetranucleotideNodes[indexTetranucleotide]);
                }

                for (var i = 0; i < data.TetranucleotideNodes[right + 1].length; ++i) {
                    index = data.TetranucleotideNodes[right + 1][i];
                    drawArrow(ctx, nodeObj.TetranucleotideNodes[indexTetranucleotide],
                        nodeObj.DinucleotideNodes[index]);
                }

                ++indexTetranucleotide;
            }
        }
    }
    
    GLOBAL.resetGraph = function resetGraph() {
        var c = document.getElementById("mycanvas");
        var ctx = c.getContext("2d");
        var width = c.width;
        var height = c.height;

        ctx.clearRect(0, 0, width, height);
    }
 });