/* global optoins */
/* global $ */
var redraw = function () { };

$(function () {
    var history = [];

    function parseObject(data, w, h) {
        var heightIdx = -1;
        var BASESES = data.Nucleotide;
        var isBaseActive = []
        var nodes = {
            TetranucleotideNodes: [],
            DinucleotideNodes: []
        };

        h -= 2 * (optoins.radius + 4)

        for (var i = 0; i < BASESES.length; ++i) {
            if (data.TetranucleotideNodes[i * 2].length +
                data.TetranucleotideNodes[i * 2 + 1].length > 0) {
                isBaseActive[i] = true
                heightIdx++
            } else {
                isBaseActive[i] = false
            }
        }

        var height = h / 4;
        var startHeight = (h - (height * heightIdx)) / 2 + optoins.radius + 4

        for (var l = 0, i = 0; l < BASESES.length; ++l) {
            if (isBaseActive[l]) {
                nodes.TetranucleotideNodes.push(nodeFactory(BASESES[l], w / 2, startHeight + height * i));
                i++;
            }
        }

        var seperators = [0]
        if (data.DinucleotideNodes.length % 2 !== 0) {
            seperators[1] = (data.DinucleotideNodes.length + 1) / 2
            seperators[2] = (data.DinucleotideNodes.length - 1) / 2
        } else {
            seperators[1] = seperators[2] = data.DinucleotideNodes.length / 2
        }




        for (var i = 1; i < seperators.length; ++i) {
            var tLength = seperators[i - 1] + seperators[i];
            if (seperators[i] > 4) {
                height = h / (seperators[i] - 1);
                startHeight = (h - (height * (seperators[i] - 1))) / 2 + optoins.radius + 4
            } else {
                height = h / seperators[i];
                startHeight = (h - (height * (seperators[i] - 1))) / 2 + optoins.radius + 4
            }
            var left = optoins.radius + 4;
            if (i == 2) {
                left = w - left
            }

            for (var l = seperators[i - 1], s = 0; l < tLength; ++l) {
                nodes.DinucleotideNodes.push(nodeFactory(data.DinucleotideNodes[l], left, height * (s) + startHeight));
                ++s;
            }
        }

        return nodes;
    }

    function nodeFactory(text, x, y) {
        return {
            text: text,
            center: {
                y: y,
                x: x
            }
        }
    }

    function drawArrow(ctx, nodeA, nodeB) {
        var fromx = nodeA.center.x;
        var fromy = nodeA.center.y;
        var tox = nodeB.center.x;
        var toy = nodeB.center.y;
        var headlen = (optoins.radius*6)/(17);   // length of head in pixels
        var angle = Math.atan2(toy - fromy, tox - fromx);

        fromx += (5 + optoins.radius) * Math.cos(angle);
        fromy += (5 + optoins.radius) * Math.sin(angle);
        tox -= (5 + optoins.radius) * Math.cos(angle); //newPoints[2]
        toy -= (5 + optoins.radius) * Math.sin(angle);

        ctx.save();
        ctx.beginPath();
        ctx.strokeStyle = optoins.arrowColor;
        ctx.moveTo(fromx, fromy);
        ctx.lineTo(tox, toy);
        ctx.lineTo(tox - headlen * Math.cos(angle - (Math.PI * optoins.radius)/ (17*10)), toy - headlen * Math.sin(angle - (Math.PI * optoins.radius)/ (6*17)));
        ctx.lineTo(tox - headlen * Math.cos(angle + (Math.PI * optoins.radius)/ (17*10)), toy - headlen * Math.sin(angle + (Math.PI * optoins.radius)/ (6*17)));
        ctx.lineTo(tox, toy);
        ctx.stroke();
        ctx.restore();
    }

    function drawGraph(nodeObj, data) {
        var c = document.getElementById("mycanvas");
        var ctx = c.getContext("2d");
        var width = optoins.width;
        var height = optoins.height;
        ctx.fillStyle = optoins.backColor;
        ctx.fillRect(0, 0, width, height)
        ctx.font = optoins.fontSize + "px Arial";
        var nodes = nodeObj.DinucleotideNodes.concat(nodeObj.TetranucleotideNodes)

        for (var i = 0; i < nodes.length; ++i) {
            ctx.fillStyle = optoins.nodeColor;
            ctx.strokeStyle = optoins.nodeBorderColor
            ctx.beginPath();
            ctx.arc(nodes[i].center.x, nodes[i].center.y, optoins.radius, 0, 2 * Math.PI);
            ctx.stroke();
            ctx.arc(nodes[i].center.x, nodes[i].center.y, optoins.radius, 0, 2 * Math.PI);
            ctx.fill();

            ctx.fillStyle = optoins.fontColor;

            var pos = ctx.measureText(nodes[i].text).width / 2;
            ctx.fillText(nodes[i].text, nodes[i].center.x - pos, ctx.fillText(nodes[i].text, nodes[i].center.x - pos, nodes[i].center.y + optoins.fontSize * 0.35));
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

    function resetGraph() {
        var c = document.getElementById("mycanvas");
        var ctx = c.getContext("2d");
        var width = c.width;
        var height = c.height;

        ctx.clearRect(0, 0, width, height);
    }

    function sendNewCodon(val, cb) {
        sendListUpdate("/newgraph", val, cb)
    }

    function sendRemoveCodon(val, cb) {
        sendListUpdate("/removecodon", val, cb)
    }

    function sendListUpdate(path, val, cb) {
        $.post(path, {
            "list": val
        }, function (data) {
            if (data !== 'Error') {
                console.log(data);
                data = JSON.parse(data);
                if (cb) {
                    cb(data)
                }

                $('#header-codon-list').text(data.List.join(', '));
                setInfo(data)
                redraw = function () {
                    resetGraph();
                    $('.canvas-container').width(optoins.width);
                    $('.canvas-container').height(optoins.height);
                    $('#mycanvas').attr('height', optoins.height);
                    $('#mycanvas').attr('width', optoins.width);
                    var nodes = parseObject(data, optoins.width, optoins.height);
                    drawGraph(nodes, data);
                };

                redraw();
            }
        });
    }

    sendNewCodon('', function (data) {
        for (var i = 0; i < data.List.length; ++i) {
            codonClick($('.' + data.List[i]));
        }
    });

    function setInfo(data) {
        if (data.CyclingIndex === 0) {
            $('#cycling-info').text("Cycling Code")
        } else {
            $('#cycling-info').text("Not Cycling Code (" + data.CyclingIndex + ")")
        }

        if (data.SelfComplementary) {
            $('#complementary-info').text("Self complementary")
        } else if (data.StrongNotSelfComplementary) {
            $('#complementary-info').text("Strong NOT self complementary")
        } else {
            $('#complementary-info').text("")
        }

        if (data.PropertyOne) {
            $('#pone-info').text("Graph has property I")
        } else {
            $('#pone-info').text("")
        }

        if (data.PropertyTwo) {
            $('#ptow-info').text("Graph has property II")
        } else {
            $('#ptow-info').text("")
        }
    }

    function codonClick($this) {
        $this.addClass('selected')
        var val = $this.text();
        history.push(val);
        var classText = $this.attr('class')
        $this.attr('class', 'temp-selected');
        var classes = classText.split(' ');
        for (var i = 0; i < classes.length; ++i) {
            if (classes[i].indexOf('class') === 0) {
                $('.' + classes[i] + ':not(no)').addClass('no')
            }
        }

        $this.attr('class', classText);

        sendNewCodon(val);
    }

    $('.codon-table td').click(function () {
        var $this = $(this);
        if ($this.hasClass('selected')) {
            history[history.length - 1] = $this.text();
            totalUndo();
        } else {
            codonClick($this);
        }
    });

    function totalReset() {
        $.post("/reset", {}, function (data) {
            if (data !== 'Error') {
                $('#header-codon-list').text("");
                resetGraph();
                resteTable();
                $('.info-container p').text("");
            }
        });
    }

    $('#reset-list').click(totalReset);

    function totalUndo() {
        if (history.length <= 1) {
            totalReset();
            return;
        }

        sendRemoveCodon(history.pop(), function (data) {
            resteTable();
            for (var i = 0; i < data.List.length; ++i) {
                codonClick($('.' + data.List[i]));
            }
        });
    }

    $('#undo-list').click(totalUndo);

    function resteTable() {
        history = [];
        $('.codon-table td.no').removeClass('no')
        $('.codon-table td.selected').removeClass('selected')
    }

});