/* global GLOBAL */
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

        h -= 2 * (GLOBAL.options.radius + 4)

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
        var minHeight = (h - (height * heightIdx)) / 2 + GLOBAL.options.radius + 4,
            startHeight = minHeight;

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
                startHeight = (h - (height * (seperators[i] - 1))) / 2 + GLOBAL.options.radius + 4
            } else {
                height = h / seperators[i];
                startHeight = (h - (height * (seperators[i] - 1))) / 2 + GLOBAL.options.radius + 4
            }
            var left = GLOBAL.options.radius + 4;
            if (i == 2) {
                left = w - left
            }

            minHeight = Math.min(minHeight, startHeight);

            for (var l = seperators[i - 1], s = 0; l < tLength; ++l) {
                nodes.DinucleotideNodes.push(nodeFactory(data.DinucleotideNodes[l], left, height * (s) + startHeight));
                ++s;
            }
        }

        return [nodes, minHeight];
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

    function sendNewCodon(val, cb) {
        sendListUpdate("/newgraph", {
            "list": val
        }, cb)
    }

    function sendRemoveCodon(val, cb) {
        sendListUpdate("/removecodon", {
            "list": val
        }, cb)
    }

    function sendNewList(val, cb) {
        sendListUpdate("/newlist", {
            "list": val
        }, cb)
    }

    function sendPermutation(val, cb) {
        sendListUpdate("/permutate", {
            "rule": val
        }, cb)
    }

    function sendShiftLeft(cb) {
        sendListUpdate("/check/shift", {}, cb)
    }

    function sendListUpdate(path, val, cb) {
        $.post(path, val, function (data) {
            if (data === 'Error') {
                if (cb) {
                    cb(data)
                }
                return;
            }

            var resizeHeight = 0,
                nodeInfo, resetData;
            if (data === 'Empty') {
                nodeInfo = false;


            } else {
                nodeInfo = true;
                data = JSON.parse(data);
                console.log(data);
                resetData = function () {
                    nodeInfo = parseObject(data, GLOBAL.options.width, GLOBAL.options.height);
                    resizeHeight = Math.max((nodeInfo[1] - 4 - GLOBAL.options.radius) * 2, 0)
                    if (cb) {
                        cb(data)
                    }

                    $('#header-codon-list').text(data.List.join(', '));
                    setInfo(data)
                }
            }

            GLOBAL.redraw = function () {
                if (nodeInfo) {
                    resetData()
                }

                GLOBAL.resetGraph();
                var $canvasContainer = $('.canvas-container')
                $canvasContainer.width(GLOBAL.options.width);
                $canvasContainer.height(GLOBAL.options.height);
                var $canvas = $('#mycanvas');
                $canvas.attr('height', GLOBAL.options.height - resizeHeight);
                $canvas.height(GLOBAL.options.height - resizeHeight);
                $canvas.css('marginTop', Math.floor(resizeHeight / 2) + 'px');
                $canvas.attr('width', GLOBAL.options.width);
                var centerHeight = $('.center-window').height();
                $('.content-conteiner').height(centerHeight);

                if (nodeInfo) {
                    GLOBAL.drawGraph(nodeInfo, data);
                }
            };

            GLOBAL.redraw();

        });
    }

    sendNewCodon('', function (data) {
        for (var i = 0; i < data.List.length; ++i) {
            codonClick($('.' + data.List[i]), true);
        }
    });

    function setInfo(data) {
        $('#code-Length').text("Length: " + data.List.length)
        $('#node-Length').text("Dinucleotide Nodes: " + data.DinucleotideNodes.length)


        if (data.MaxPath === 0) {
            $('#path-info').text("-")
        } else {
            $('#path-info').text("Max path length: (" + data.MaxPath + ")")
        }
        
        if (data.CyclingIndex === 0) {
            $('#cycling-info').text("Cycling-Code")
        } else {
            $('#cycling-info').text("Not cycling-code (" + data.CyclingIndex + ")")
        }

        if (data.SelfComplementary) {
            $('#complementary-info').text("Self-complementary");
        } else if (data.StrongNotSelfComplementary) {
            $('#complementary-info').text("Strong not self-complementary");
        } else {
            $('#complementary-info').text("");
        }

        if (data.PropertyOne) {
            $('#pone-info').text("Graph has property I");
        } else {
            $('#pone-info').text("");
        }

        if (data.PropertyTwo) {
            $('#ptow-info').text("Graph has property II");
        } else {
            $('#ptow-info').text("");
        }

        $('#full-cycling-info').text("");
    }

    function codonClick($this, isLocal) {
        $this.addClass('selected')
        var val = $this.text();
        history.push(val);
        var classText = [];
        $this.each(function () {
            classText.push($(this).attr('class'))
        });

        var classes = classText[0].split(' ');
        for (var i = 0; i < classes.length; ++i) {
            if (classes[i].indexOf('class') === 0) {
                $('.' + classes[i] + ':not(no)').addClass('no');
            }
        }
        var classIndex = 0;
        $this.each(function () {
            $(this).attr('class', classText[classIndex++]);
        });

        $('.' + val).addClass('selected').removeClass('no');
        if (!isLocal) {
            sendNewCodon(val);
        }
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

    $('#button-permutation-rule').click(function () {
        $('#header-codon-list').text("");
        $('.info-container p').text("");
        GLOBAL.resetGraph();
        resteTable();
        sendPermutation($('#select-permutation-rule').val(), function (data) {
            for (var i = 0; i < data.List.length; ++i) {
                codonClick($('.' + data.List[i]), true);
            }
        })
    });

    $('#button-shift').click(function () {
        sendShiftLeft(function (data) {
            $('#header-codon-list').text("");
            $('.info-container p').text("");
            GLOBAL.resetGraph();
            resteTable();
            if (data !== "Error") {
                for (var i = 0; i < data.List.length; ++i) {
                    codonClick($('.' + data.List[i]), true);
                }
            }
        });
    });

    function totalReset() {
        $.post("/reset", {}, function (data) {
            if (data !== 'Error') {
                $('#header-codon-list').text("");
                GLOBAL.resetGraph();
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
                codonClick($('.' + data.List[i]), true);
            }
        });
    }

    // $('#undo-list').click(totalUndo);

    function resteTable() {
        history = [];
        $('.codon-table td.no').removeClass('no')
        $('.codon-table td.selected').removeClass('selected')
    }

    $('#send-list').click(function () {
        var value = $('#codonList').val().toUpperCase();
        $('#header-codon-list').text("");
        GLOBAL.resetGraph();
        resteTable();
        $('.info-container p').text("");
        sendNewList(value, function (data) {
            for (var i = 0; i < data.List.length; ++i) {
                codonClick($('.' + data.List[i]), true);
            }
        });

    });
});