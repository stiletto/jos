
function JosVis(canvas, form) {
    this.canvas = canvas;
    this.form = form;
    this.inputNodes = form.elements["nodes"];
    this.inputCopies = form.elements["copies"];

    function normalize_rgb_value(color, m) {
        color = Math.floor((color + m) * 255);
        if (color < 0) {
            color = 0;
        }
        return color;
    }
    function rgbToHex(r, g, b) {
        return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
    }
    this.HslToRgb = function (h,s,l) {
        var c = (1 - Math.abs(2*l - 1)) * s,
            x = c * ( 1 - Math.abs((h / 60 ) % 2 - 1 )),
            m = l - c/ 2,
            r, g, b;

        if (h < 60) {
            r = c;
            g = x;
            b = 0;
        }
        else if (h < 120) {
            r = x;
            g = c;
            b = 0;
        }
        else if (h < 180) {
            r = 0;
            g = c;
            b = x;
        }
        else if (h < 240) {
            r = 0;
            g = x;
            b = c;
        }
        else if (h < 300) {
            r = x;
            g = 0;
            b = c;
        }
        else {
            r = c;
            g = 0;
            b = x;
        }

		r = normalize_rgb_value(r, m);
		g = normalize_rgb_value(g, m);
		b = normalize_rgb_value(b, m);

		return rgbToHex(r,g,b);
	}
	var colors = [];
	var colorCount = 64;
	for (var i=0; i<colorCount; i++) {
		colors.push(this.HslToRgb(Math.trunc(i%8)*360/8, 1.0-0.4*Math.trunc(i/8)%8, 0.5));
	}

    this.update = function () {
        var nodes = Number(this.inputNodes.value);
        var copies = Number(this.inputCopies.value);
        var ctx = canvas.getContext('2d');
        ctx.clearRect(0, 0, canvas.width, canvas.height);
        ctx.save();
        ctx.translate(canvas.width/2,canvas.height/2);
        ctx.rotate(-90*Math.PI/180);
        ctx.save();
        var nodeSec = 2*Math.PI/nodes;
        var maxR = canvas.height/4 - 16;
        for (var copy=0; copy<copies; copy++) {
			var copyR = maxR + copy*32;
			for (var node=0; node<nodes; node++) {
				for (var step=0; step<2; step++) {
					ctx.beginPath();
					var ars = node*nodeSec + copy/copies*2*Math.PI;
					var are = (node+1)*nodeSec + copy/copies*2*Math.PI;
					
					ctx.arc(0,0, copyR, ars, are);
					ctx.arc(0,0, copyR-16, are, ars, true);
					if (step==0) {
						ctx.fillStyle = colors[node%colors.length];
						ctx.fill();
					} else {
						ctx.stroke();
					}
				}
			}
		}
		ctx.beginPath()
		ctx.moveTo(0,0);
		ctx.lineTo(canvas.height/2-16,0);
		ctx.stroke();
        ctx.restore();
        ctx.restore();
    };
    var jv = this;
    form.addEventListener("submit", function(evt) {
        jv.update();
        evt.preventDefault();
    }, true);
}

var jv = new JosVis(document.getElementById("cv"), document.getElementById('inform'));
jv.canvas.width = window.innerWidth;
jv.canvas.height = 700;
jv.update();
