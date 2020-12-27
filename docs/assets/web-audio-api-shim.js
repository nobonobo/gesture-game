(() => {
  var __commonJS = (callback, module) => () => {
    if (!module) {
      module = {exports: {}};
      callback(module.exports, module);
    }
    return module.exports;
  };

  // node_modules/@mohayonao/web-audio-api-shim/lib/AnalyserNode.js
  var require_AnalyserNode = __commonJS((exports) => {
    "use strict";
    Object.defineProperty(exports, "__esModule", {
      value: true
    });
    exports.install = install;
    var AnalyserNode = global.AnalyserNode;
    function installGetFloatTimeDomainData() {
      if (AnalyserNode.prototype.hasOwnProperty("getFloatTimeDomainData")) {
        return;
      }
      var uint8 = new Uint8Array(2048);
      AnalyserNode.prototype.getFloatTimeDomainData = function(array) {
        this.getByteTimeDomainData(uint8);
        for (var i = 0, imax = array.length; i < imax; i++) {
          array[i] = (uint8[i] - 128) * 78125e-7;
        }
      };
    }
    function install() {
      installGetFloatTimeDomainData();
    }
  });

  // node_modules/@mohayonao/web-audio-api-shim/lib/AudioBuffer.js
  var require_AudioBuffer = __commonJS((exports) => {
    "use strict";
    Object.defineProperty(exports, "__esModule", {
      value: true
    });
    exports.install = install;
    var AudioBuffer = global.AudioBuffer;
    function installCopyFromChannel() {
      if (AudioBuffer.prototype.hasOwnProperty("copyFromChannel")) {
        return;
      }
      AudioBuffer.prototype.copyFromChannel = function(destination, channelNumber, startInChannel) {
        var source = this.getChannelData(channelNumber | 0).subarray(startInChannel | 0);
        destination.set(source.subarray(0, Math.min(source.length, destination.length)));
      };
    }
    function installCopyToChannel() {
      if (AudioBuffer.prototype.hasOwnProperty("copyToChannel")) {
        return;
      }
      AudioBuffer.prototype.copyToChannel = function(source, channelNumber, startInChannel) {
        var clipped = source.subarray(0, Math.min(source.length, this.length - (startInChannel | 0)));
        this.getChannelData(channelNumber | 0).set(clipped, startInChannel | 0);
      };
    }
    function install() {
      installCopyFromChannel();
      installCopyToChannel();
    }
  });

  // node_modules/@mohayonao/web-audio-api-shim/lib/AudioNode.js
  var require_AudioNode = __commonJS((exports) => {
    "use strict";
    Object.defineProperty(exports, "__esModule", {
      value: true
    });
    exports.install = install;
    var OfflineAudioContext = global.OfflineAudioContext;
    var AudioNode = global.AudioNode;
    var connect = AudioNode.prototype.connect;
    var disconnect = AudioNode.prototype.disconnect;
    function match(args, connection) {
      for (var i = 0, imax = args.length; i < imax; i++) {
        if (args[i] !== connection[i]) {
          return false;
        }
      }
      return true;
    }
    function disconnectAll(node) {
      for (var ch = 0, chmax = node.numberOfOutputs; ch < chmax; ch++) {
        disconnect.call(node, ch);
      }
      node._shim$connections = [];
    }
    function disconnectChannel(node, channel) {
      disconnect.call(node, channel);
      node._shim$connections = node._shim$connections.filter(function(connection) {
        return connection[1] !== channel;
      });
    }
    function disconnectSelect(node, args) {
      var remain = [];
      var hasDestination = false;
      node._shim$connections.forEach(function(connection) {
        hasDestination = hasDestination || args[0] === connection[0];
        if (!match(args, connection)) {
          remain.push(connection);
        }
      });
      if (!hasDestination) {
        throw new Error("Failed to execute 'disconnect' on 'AudioNode': the given destination is not connected.");
      }
      disconnectAll(node);
      remain.forEach(function(connection) {
        connect.call(node, connection[0], connection[1], connection[2]);
      });
      node._shim$connections = remain;
    }
    function installDisconnect() {
      var audioContext = new OfflineAudioContext(1, 1, 44100);
      var isSelectiveDisconnection = false;
      try {
        audioContext.createGain().disconnect(audioContext.destination);
      } catch (e) {
        isSelectiveDisconnection = true;
      }
      if (isSelectiveDisconnection) {
        return;
      }
      AudioNode.prototype.disconnect = function() {
        this._shim$connections = this._shim$connections || [];
        for (var _len = arguments.length, args = Array(_len), _key = 0; _key < _len; _key++) {
          args[_key] = arguments[_key];
        }
        if (args.length === 0) {
          disconnectAll(this);
        } else if (args.length === 1 && typeof args[0] === "number") {
          disconnectChannel(this, args[0]);
        } else {
          disconnectSelect(this, args);
        }
      };
      AudioNode.prototype.disconnect.original = disconnect;
      AudioNode.prototype.connect = function(destination) {
        var output = arguments[1] === void 0 ? 0 : arguments[1];
        var input = arguments[2] === void 0 ? 0 : arguments[2];
        var _input = void 0;
        this._shim$connections = this._shim$connections || [];
        if (destination instanceof AudioNode) {
          connect.call(this, destination, output, input);
          _input = input;
        } else {
          connect.call(this, destination, output);
          _input = 0;
        }
        this._shim$connections.push([destination, output, _input]);
      };
      AudioNode.prototype.connect.original = connect;
    }
    function install(stage) {
      if (stage !== 0) {
        installDisconnect();
      }
    }
  });

  // node_modules/stereo-panner-node/lib/curve.js
  var require_curve = __commonJS((exports, module) => {
    var WS_CURVE_SIZE = 4096;
    var curveL = new Float32Array(WS_CURVE_SIZE);
    var curveR = new Float32Array(WS_CURVE_SIZE);
    (function() {
      var i;
      for (i = 0; i < WS_CURVE_SIZE; i++) {
        curveL[i] = Math.cos(i / WS_CURVE_SIZE * Math.PI * 0.5);
        curveR[i] = Math.sin(i / WS_CURVE_SIZE * Math.PI * 0.5);
      }
    })();
    module.exports = {
      L: curveL,
      R: curveR
    };
  });

  // node_modules/stereo-panner-node/lib/stereo-panner-impl.js
  var require_stereo_panner_impl = __commonJS((exports, module) => {
    var curve = require_curve();
    function StereoPannerImpl(audioContext) {
      this.audioContext = audioContext;
      this.inlet = audioContext.createChannelSplitter(2);
      this._pan = audioContext.createGain();
      this.pan = this._pan.gain;
      this._wsL = audioContext.createWaveShaper();
      this._wsR = audioContext.createWaveShaper();
      this._L = audioContext.createGain();
      this._R = audioContext.createGain();
      this.outlet = audioContext.createChannelMerger(2);
      this.inlet.channelCount = 2;
      this.inlet.channelCountMode = "explicit";
      this._pan.gain.value = 0;
      this._wsL.curve = curve.L;
      this._wsR.curve = curve.R;
      this._L.gain.value = 0;
      this._R.gain.value = 0;
      this.inlet.connect(this._L, 0);
      this.inlet.connect(this._R, 1);
      this._L.connect(this.outlet, 0, 0);
      this._R.connect(this.outlet, 0, 1);
      this._pan.connect(this._wsL);
      this._pan.connect(this._wsR);
      this._wsL.connect(this._L.gain);
      this._wsR.connect(this._R.gain);
      this._isConnected = false;
      this._dc1buffer = null;
      this._dc1 = null;
    }
    StereoPannerImpl.prototype.connect = function(destination) {
      var audioContext = this.audioContext;
      if (!this._isConnected) {
        this._isConnected = true;
        this._dc1buffer = audioContext.createBuffer(1, 2, audioContext.sampleRate);
        this._dc1buffer.getChannelData(0).set([1, 1]);
        this._dc1 = audioContext.createBufferSource();
        this._dc1.buffer = this._dc1buffer;
        this._dc1.loop = true;
        this._dc1.start(audioContext.currentTime);
        this._dc1.connect(this._pan);
      }
      global.AudioNode.prototype.connect.call(this.outlet, destination);
    };
    StereoPannerImpl.prototype.disconnect = function() {
      var audioContext = this.audioContext;
      if (this._isConnected) {
        this._isConnected = false;
        this._dc1.stop(audioContext.currentTime);
        this._dc1.disconnect();
        this._dc1 = null;
        this._dc1buffer = null;
      }
      global.AudioNode.prototype.disconnect.call(this.outlet);
    };
    module.exports = StereoPannerImpl;
  });

  // node_modules/stereo-panner-node/lib/stereo-panner-node.js
  var require_stereo_panner_node = __commonJS((exports, module) => {
    var StereoPannerImpl = require_stereo_panner_impl();
    var AudioContext = global.AudioContext || global.webkitAudioContext;
    function StereoPanner(audioContext) {
      var impl = new StereoPannerImpl(audioContext);
      Object.defineProperties(impl.inlet, {
        pan: {
          value: impl.pan,
          enumerable: true
        },
        connect: {
          value: function(node) {
            return impl.connect(node);
          }
        },
        disconnect: {
          value: function() {
            return impl.disconnect();
          }
        }
      });
      return impl.inlet;
    }
    StereoPanner.polyfill = function() {
      if (!AudioContext || AudioContext.prototype.hasOwnProperty("createStereoPanner")) {
        return;
      }
      AudioContext.prototype.createStereoPanner = function() {
        return new StereoPanner(this);
      };
    };
    module.exports = StereoPanner;
  });

  // node_modules/@mohayonao/web-audio-api-shim/lib/AudioContext.js
  var require_AudioContext = __commonJS((exports) => {
    "use strict";
    Object.defineProperty(exports, "__esModule", {
      value: true
    });
    var _createClass = function() {
      function defineProperties(target, props) {
        for (var i = 0; i < props.length; i++) {
          var descriptor = props[i];
          descriptor.enumerable = descriptor.enumerable || false;
          descriptor.configurable = true;
          if ("value" in descriptor)
            descriptor.writable = true;
          Object.defineProperty(target, descriptor.key, descriptor);
        }
      }
      return function(Constructor, protoProps, staticProps) {
        if (protoProps)
          defineProperties(Constructor.prototype, protoProps);
        if (staticProps)
          defineProperties(Constructor, staticProps);
        return Constructor;
      };
    }();
    var _get = function get(_x, _x2, _x3) {
      var _again = true;
      _function:
        while (_again) {
          var object = _x, property = _x2, receiver = _x3;
          desc = parent = getter = void 0;
          _again = false;
          if (object === null)
            object = Function.prototype;
          var desc = Object.getOwnPropertyDescriptor(object, property);
          if (desc === void 0) {
            var parent = Object.getPrototypeOf(object);
            if (parent === null) {
              return void 0;
            } else {
              _x = parent;
              _x2 = property;
              _x3 = receiver;
              _again = true;
              continue _function;
            }
          } else if ("value" in desc) {
            return desc.value;
          } else {
            var getter = desc.get;
            if (getter === void 0) {
              return void 0;
            }
            return getter.call(receiver);
          }
        }
    };
    exports.install = install;
    function _classCallCheck(instance, Constructor) {
      if (!(instance instanceof Constructor)) {
        throw new TypeError("Cannot call a class as a function");
      }
    }
    function _inherits(subClass, superClass) {
      if (typeof superClass !== "function" && superClass !== null) {
        throw new TypeError("Super expression must either be null or a function, not " + typeof superClass);
      }
      subClass.prototype = Object.create(superClass && superClass.prototype, {constructor: {value: subClass, enumerable: false, writable: true, configurable: true}});
      if (superClass)
        subClass.__proto__ = superClass;
    }
    var OriginalAudioContext = global.AudioContext;
    var OriginalOfflineAudioContext = global.OfflineAudioContext;
    var AudioNode = global.AudioNode;
    var EventTarget = global.EventTarget || global.Object.constructor;
    function nop() {
    }
    function inherits(ctor, superCtor) {
      ctor.prototype = Object.create(superCtor.prototype, {
        constructor: {value: ctor, enumerable: false, writable: true, configurable: true}
      });
    }
    function replaceAudioContext() {
      if (global.AudioContext !== OriginalAudioContext) {
        return;
      }
      function BaseAudioContext(audioContext) {
        this._ = {};
        this._.audioContext = audioContext;
        this._.destination = audioContext.destination;
        this._.state = "";
        this._.currentTime = 0;
        this._.sampleRate = audioContext.sampleRate;
        this._.onstatechange = null;
      }
      inherits(BaseAudioContext, EventTarget);
      Object.defineProperties(BaseAudioContext.prototype, {
        destination: {
          get: function get() {
            return this._.destination;
          }
        },
        sampleRate: {
          get: function get() {
            return this._.sampleRate;
          }
        },
        currentTime: {
          get: function get() {
            return this._.currentTime || this._.audioContext.currentTime;
          }
        },
        listener: {
          get: function get() {
            return this._.audioContext.listener;
          }
        },
        state: {
          get: function get() {
            return this._.state;
          }
        },
        onstatechange: {
          set: function set(fn) {
            if (typeof fn === "function") {
              this._.onstatechange = fn;
            }
          },
          get: function get() {
            return this._.onstatechange;
          }
        }
      });
      var AudioContext = function(_BaseAudioContext) {
        function AudioContext2() {
          _classCallCheck(this, AudioContext2);
          _get(Object.getPrototypeOf(AudioContext2.prototype), "constructor", this).call(this, new OriginalAudioContext());
          this._.state = "running";
          if (!OriginalAudioContext.prototype.hasOwnProperty("suspend")) {
            this._.destination = this._.audioContext.createGain();
            this._.destination.connect(this._.audioContext.destination);
            this._.destination.connect = function() {
              this._.audioContext.destination.connect.apply(this._.audioContext.destination, arguments);
            };
            this._.destination.disconnect = function() {
              this._.audioContext.destination.connect.apply(this._.audioContext.destination, arguments);
            };
            this._.destination.channelCountMode = "explicit";
          }
        }
        _inherits(AudioContext2, _BaseAudioContext);
        return AudioContext2;
      }(BaseAudioContext);
      AudioContext.prototype.suspend = function() {
        var _this = this;
        if (this._.state === "closed") {
          return Promise.reject(new Error("cannot suspend a closed AudioContext"));
        }
        function changeState() {
          this._.state = "suspended";
          this._.currentTime = this._.audioContext.currentTime;
        }
        var promise = void 0;
        if (typeof this._.audioContext === "function") {
          promise = this._.audioContext.suspend();
          promise.then(function() {
            changeState.call(_this);
          });
        } else {
          AudioNode.prototype.disconnect.call(this._.destination);
          promise = Promise.resolve();
          promise.then(function() {
            changeState.call(_this);
            var e = new global.Event("statechange");
            if (typeof _this._.onstatechange === "function") {
              _this._.onstatechange(e);
            }
            _this.dispatchEvent(e);
          });
        }
        return promise;
      };
      AudioContext.prototype.resume = function() {
        var _this2 = this;
        if (this._.state === "closed") {
          return Promise.reject(new Error("cannot resume a closed AudioContext"));
        }
        function changeState() {
          this._.state = "running";
          this._.currentTime = 0;
        }
        var promise = void 0;
        if (typeof this._.audioContext.resume === "function") {
          promise = this._.audioContext.resume();
          promise.then(function() {
            changeState.call(_this2);
          });
        } else {
          AudioNode.prototype.connect.call(this._.destination, this._.audioContext.destination);
          promise = Promise.resolve();
          promise.then(function() {
            changeState.call(_this2);
            var e = new global.Event("statechange");
            if (typeof _this2._.onstatechange === "function") {
              _this2._.onstatechange(e);
            }
            _this2.dispatchEvent(e);
          });
        }
        return promise;
      };
      AudioContext.prototype.close = function() {
        var _this3 = this;
        if (this._.state === "closed") {
          return Promise.reject(new Error("Cannot close a context that is being closed or has already been closed."));
        }
        function changeState() {
          this._.state = "closed";
          this._.currentTime = Infinity;
          this._.sampleRate = 0;
        }
        var promise = void 0;
        if (typeof this._.audioContext.close === "function") {
          promise = this._.audioContext.close();
          promise.then(function() {
            changeState.call(_this3);
          });
        } else {
          if (typeof this._.audioContext.suspend === "function") {
            this._.audioContext.suspend();
          } else {
            AudioNode.prototype.disconnect.call(this._.destination);
          }
          promise = Promise.resolve();
          promise.then(function() {
            changeState.call(_this3);
            var e = new global.Event("statechange");
            if (typeof _this3._.onstatechange === "function") {
              _this3._.onstatechange(e);
            }
            _this3.dispatchEvent(e);
          });
        }
        return promise;
      };
      ["addEventListener", "removeEventListener", "dispatchEvent", "createBuffer"].forEach(function(methodName) {
        AudioContext.prototype[methodName] = function() {
          return this._.audioContext[methodName].apply(this._.audioContext, arguments);
        };
      });
      ["decodeAudioData", "createBufferSource", "createMediaElementSource", "createMediaStreamSource", "createMediaStreamDestination", "createAudioWorker", "createScriptProcessor", "createAnalyser", "createGain", "createDelay", "createBiquadFilter", "createWaveShaper", "createPanner", "createStereoPanner", "createConvolver", "createChannelSplitter", "createChannelMerger", "createDynamicsCompressor", "createOscillator", "createPeriodicWave"].forEach(function(methodName) {
        AudioContext.prototype[methodName] = function() {
          if (this._.state === "closed") {
            throw new Error("Failed to execute '" + methodName + "' on 'AudioContext': AudioContext has been closed");
          }
          return this._.audioContext[methodName].apply(this._.audioContext, arguments);
        };
      });
      var OfflineAudioContext = function(_BaseAudioContext2) {
        function OfflineAudioContext2(numberOfChannels, length, sampleRate) {
          _classCallCheck(this, OfflineAudioContext2);
          _get(Object.getPrototypeOf(OfflineAudioContext2.prototype), "constructor", this).call(this, new OriginalOfflineAudioContext(numberOfChannels, length, sampleRate));
          this._.state = "suspended";
        }
        _inherits(OfflineAudioContext2, _BaseAudioContext2);
        _createClass(OfflineAudioContext2, [{
          key: "oncomplete",
          set: function set(fn) {
            this._.audioContext.oncomplete = fn;
          },
          get: function get() {
            return this._.audioContext.oncomplete;
          }
        }]);
        return OfflineAudioContext2;
      }(BaseAudioContext);
      ["addEventListener", "removeEventListener", "dispatchEvent", "createBuffer", "decodeAudioData", "createBufferSource", "createMediaElementSource", "createMediaStreamSource", "createMediaStreamDestination", "createAudioWorker", "createScriptProcessor", "createAnalyser", "createGain", "createDelay", "createBiquadFilter", "createWaveShaper", "createPanner", "createStereoPanner", "createConvolver", "createChannelSplitter", "createChannelMerger", "createDynamicsCompressor", "createOscillator", "createPeriodicWave"].forEach(function(methodName) {
        OfflineAudioContext.prototype[methodName] = function() {
          return this._.audioContext[methodName].apply(this._.audioContext, arguments);
        };
      });
      OfflineAudioContext.prototype.startRendering = function() {
        var _this4 = this;
        if (this._.state !== "suspended") {
          return Promise.reject(new Error("cannot call startRendering more than once"));
        }
        this._.state = "running";
        var promise = this._.audioContext.startRendering();
        promise.then(function() {
          _this4._.state = "closed";
          var e = new global.Event("statechange");
          if (typeof _this4._.onstatechange === "function") {
            _this4._.onstatechange(e);
          }
          _this4.dispatchEvent(e);
        });
        return promise;
      };
      OfflineAudioContext.prototype.suspend = function() {
        if (typeof this._.audioContext.suspend === "function") {
          return this._.audioContext.suspend();
        }
        return Promise.reject(new Error("cannot suspend an OfflineAudioContext"));
      };
      OfflineAudioContext.prototype.resume = function() {
        if (typeof this._.audioContext.resume === "function") {
          return this._.audioContext.resume();
        }
        return Promise.reject(new Error("cannot resume an OfflineAudioContext"));
      };
      OfflineAudioContext.prototype.close = function() {
        if (typeof this._.audioContext.close === "function") {
          return this._.audioContext.close();
        }
        return Promise.reject(new Error("cannot close an OfflineAudioContext"));
      };
      global.AudioContext = AudioContext;
      global.OfflineAudioContext = OfflineAudioContext;
    }
    function installCreateAudioWorker() {
    }
    function installCreateStereoPanner() {
      if (OriginalAudioContext.prototype.hasOwnProperty("createStereoPanner")) {
        return;
      }
      var StereoPannerNode = require_stereo_panner_node();
      OriginalAudioContext.prototype.createStereoPanner = function() {
        return new StereoPannerNode(this);
      };
    }
    function installDecodeAudioData() {
      var audioContext = new OriginalOfflineAudioContext(1, 1, 44100);
      var isPromiseBased = false;
      try {
        var audioData = new Uint32Array([1179011410, 48, 1163280727, 544501094, 16, 131073, 44100, 176400, 1048580, 1635017060, 8, 0, 0, 0, 0]).buffer;
        isPromiseBased = !!audioContext.decodeAudioData(audioData, nop);
      } catch (e) {
        nop(e);
      }
      if (isPromiseBased) {
        return;
      }
      var decodeAudioData = OriginalAudioContext.prototype.decodeAudioData;
      OriginalAudioContext.prototype.decodeAudioData = function(audioData2, successCallback, errorCallback) {
        var _this5 = this;
        var promise = new Promise(function(resolve, reject) {
          return decodeAudioData.call(_this5, audioData2, resolve, reject);
        });
        promise.then(successCallback, errorCallback);
        return promise;
      };
      OriginalAudioContext.prototype.decodeAudioData.original = decodeAudioData;
    }
    function installClose() {
      if (OriginalAudioContext.prototype.hasOwnProperty("close")) {
        return;
      }
      replaceAudioContext();
    }
    function installResume() {
      if (OriginalAudioContext.prototype.hasOwnProperty("resume")) {
        return;
      }
      replaceAudioContext();
    }
    function installSuspend() {
      if (OriginalAudioContext.prototype.hasOwnProperty("suspend")) {
        return;
      }
      replaceAudioContext();
    }
    function installStartRendering() {
      var audioContext = new OriginalOfflineAudioContext(1, 1, 44100);
      var isPromiseBased = false;
      try {
        isPromiseBased = !!audioContext.startRendering();
      } catch (e) {
        nop(e);
      }
      if (isPromiseBased) {
        return;
      }
      var startRendering = OriginalOfflineAudioContext.prototype.startRendering;
      OriginalOfflineAudioContext.prototype.startRendering = function() {
        var _this6 = this;
        return new Promise(function(resolve) {
          var oncomplete = _this6.oncomplete;
          _this6.oncomplete = function(e) {
            resolve(e.renderedBuffer);
            if (typeof oncomplete === "function") {
              oncomplete.call(_this6, e);
            }
          };
          startRendering.call(_this6);
        });
      };
      OriginalOfflineAudioContext.prototype.startRendering.original = startRendering;
    }
    function install(stage) {
      installCreateAudioWorker();
      installCreateStereoPanner();
      installDecodeAudioData();
      installStartRendering();
      if (stage !== 0) {
        installClose();
        installResume();
        installSuspend();
      }
    }
  });

  // node_modules/@mohayonao/web-audio-api-shim/lib/install.js
  var require_install = __commonJS((exports, module) => {
    "use strict";
    Object.defineProperty(exports, "__esModule", {
      value: true
    });
    exports["default"] = install;
    function install() {
      var stage = arguments[0] === void 0 ? Infinity : arguments[0];
      if (!global.hasOwnProperty("AudioContext") && global.hasOwnProperty("webkitAudioContext")) {
        global.AudioContext = global.webkitAudioContext;
      }
      if (!global.hasOwnProperty("OfflineAudioContext") && global.hasOwnProperty("webkitOfflineAudioContext")) {
        global.OfflineAudioContext = global.webkitOfflineAudioContext;
      }
      if (!global.AudioContext) {
        return;
      }
      require_AnalyserNode().install(stage);
      require_AudioBuffer().install(stage);
      require_AudioNode().install(stage);
      require_AudioContext().install(stage);
    }
    module.exports = exports["default"];
  });

  // node_modules/@mohayonao/web-audio-api-shim/index.js
  var require_web_audio_api_shim = __commonJS((exports, module) => {
    module.exports = require_install()(Infinity);
  });

  // index.js
  require_web_audio_api_shim();
})();
