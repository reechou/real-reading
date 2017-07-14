function shareSucceedCallback() {
    $playStatus.share()
}
$payButton = {
    init: function () {
        $payButton.targetUrl(), $(".scan_window_close").click(function () {
            $pay.isPaying = !1, $(".scan_layer").css("display", "none"), $(".navigation_menu").css("display", "block")
        }), $(".scan_window_button").click(function () {
            "undefined" != typeof Storage && (localStorage.show_qr_layer = "show");
            var o = new Date;
            eJump(location.href + "?" + o.getTime(), !0)
        })
    }, targetUrl: function () {
        $("#payButton").on("click", function () {
            if (3 === parseInt($("#paymentType").val())) 1 === parseInt($("#canPay").val()) ? $pay.pay(!0) : location.href = $("#paymentUrl").val();
            else {
                var o = $("#coupons").val();
                console.log(o.length), o.length >= 5 ? ($(".column-pay").hide(), $(".darkScreen").show(), $(".couponPaymentBox").show(), $(".darkScreen").click(function () {
                    $(".darkScreen").hide(), $(".couponPaymentBox").hide(), $(".column-pay").show()
                })) : $pay.pay()
            }
            _czc.push(["_trackEvent", "btn_buy", "btn_click", "购买支付"])
        }), $("#columnPayButton").on("click", function () {
            location.href = $("#paymentUrl").val()
        })
    }
}, $pay = {
    isPaying: !1,
    scanUrl: "",
    hasPhone: !0,
    pay: function (o) {
        if (!$pay.isPaying) {
            if ($pay.isPaying = !0, $pay.scanUrl.length > 0) return $(".scan_window_image").attr("src", $pay.scanUrl), $(".scan_layer").css("display", "block"), void $(".navigation_menu").css("display", "none");
            var i = null;
            o && (i = {
                buy_product: 1
            }), NetWork.request("pay/get_info", i, function (o) {
                var i = $.parseJSON(o);
                if (0 == i.code) {
                    $pay.hasPhone = i.data.hasPhone;
                    var a = i.data.payConfig,
                        t = a["package"].split("=")[1];
                    if (i.data.hasOwnProperty("pay_url")) localStorage.setItem("need_reload_when_show", "1"), window.location.href = i.data.pay_url;
                    else {
                        if ("undefined" == typeof WeixinJSBridge) return void($pay.isPaying = !1);
                        WeixinJSBridge.invoke("getBrandWCPayRequest", {
                            appId: a.appId,
                            timeStamp: a.timeStamp,
                            nonceStr: a.nonceStr,
                            "package": a["package"],
                            signType: a.signType,
                            paySign: a.paySign
                        }, function (o) {
                            if ($pay.isPaying = !1, o.err_code) {
                                var i = {
                                    type: "user_pay_err_per_minute",
                                    extra: "js"
                                };
                                if (NetWork.request("log_watcher_data", i, function (o) {}), 3 == o.err_code || 2 == o.err_code) {
                                    var a = window.location.hostname.indexOf(".h5.");
                                    a === -1 ? $pay.scanUrl = "/" + window.APPID + "/pay/scan_image/" + t : $pay.scanUrl = "/pay/scan_image/" + t, $(".scan_window_image").attr("src", $pay.scanUrl), $(".scan_layer").css("display", "block"), $(".navigation_menu").css("display", "none")
                                } else alert(o.err_code + " " + o.err_msg)
                            } else if ("get_brand_wcpay_request:ok" == o.err_msg) {
                                "undefined" != typeof Storage && (localStorage.show_qr_layer = "show");
                                var e = (new Date).getTime(),
                                    n = location.href + "?" + e;
                                eJump(n, !0)
                            } else "get_brand_wcpay_request:cancel" == o.err_msg
                        })
                    }
                }
            })
        }
    }
}, $subscribeModule = {
    bodyTouchMoveListener: null,
    init: function () {
        $subscribeModule.bodyTouchMoveListener = function (o) {
            o.preventDefault()
        };
        var o = $(".pay_qr_layer").data("available"),
            i = $(".pay_qr_url").attr("src");
        "undefined" != typeof Storage && "show" == localStorage.show_qr_layer && 1 == o && i && (document.body.addEventListener("touchmove", $subscribeModule.bodyTouchMoveListener, !1), $(".pay_qr_layer").css("display", "block"), localStorage.show_qr_layer = "hide"), $(".pay_qr_close").click(function () {
            document.body.removeEventListener("touchmove", $subscribeModule.bodyTouchMoveListener), $(".pay_qr_layer").css("display", "none")
        })
    }
}, $hideDetail = {
    $hideStauts: !1,
    init: function () {
        $("#more_detail_wrapper").on("click", function () {
            $hideDetail.$hideStauts ? ($(".detail-desc").show(), $("#more_detail_button").removeClass("more_detail_icon").addClass("less_detail_icon"), $hideDetail.$hideStauts = !1) : ($(".detail-desc").hide(), $("#more_detail_button").removeClass("less_detail_icon").addClass("more_detail_icon"), $hideDetail.$hideStauts = !0), $audioDetail.$droploadMe && $audioDetail.$droploadMe.resetload()
        })
    }
}, $audioDetail = {
    $comment_bottom_face: $("#comment_bottom_face"),
    $comment_input_face: $("#comment_input_face"),
    $comment_bottom: $("#comment_bottom"),
    $comment_input_real: $("#comment_input_real"),
    $scrollArea: null,
    $maxHeight: null,
    $minHeight: null,
    $isKeyBoardShown: !1,
    $globalFontSize: !1,
    $commentsListTopHot: $("#comments_list_tophot"),
    $commentsList: $("#comments_list"),
    $isSending: null,
    $isLoading: null,
    $firtContent: !1,
    $droploadMe: null,
    $signButton: null,
    $signLayer: !1,
    $signLayerClose: !1,
    $signLayerBigImage: !1,
    $signImgUrl: null,
    $src_commentid: "",
    $src_only_for_group_post_userid: "",
    $src_userid: "",
    $src_content: "",
    $src_nickname: "",
    $scrollTop: 0,
    $currentCommentPage: 1,
    $isOrderByHot: !1,
    init: function () {
        $audioDetail.dealTopCommentsEmoji(), $audioDetail.$maxHeight = window.innerHeight, $audioDetail.$globalFontSize = document.documentElement.clientWidth / 3.75, $audioDetail.$comment_input_face.click(function () {
            $audioDetail.$scrollTop = $("body").scrollTop(), $("#comment_bottom").on("touchmove", function (o) {
                o.preventDefault()
            }), window.scrollTo(0, 0), $audioDetail.$comment_bottom_face.hide(), $audioDetail.$comment_bottom.show(), $audioDetail.$comment_input_real.focus()
        }), $audioDetail.initCommentList(), $audioDetail.loadMoreComment(null), $("#phone_cancel_button").click(function () {
            $audioDetail.$scrollTop && ($("body").scrollTop($audioDetail.$scrollTop), $audioDetail.$scrollTop = 0), $("body").off("touchmove"), $audioDetail.$comment_bottom.hide(), $audioDetail.$comment_bottom_face.show()
        }), $("#send_ok_btn").click(function () {
            if ($("body").off("touchmove"), !$audioDetail.$isSending) {
                var o = $audioDetail.$comment_input_real.val();
                if ("" == o.replace(/ /g, "") || "" == o || null == o) $audioDetail.$comment_input_real.val(""), $audioDetail.$comment_input_face.html("发评论"), $audioDetail.$comment_input_real.attr("placeholder", "发评论"), $audioDetail.clearReplyInfos();
                else {
                    o = $emojiModule.emojiMapping(o), o = $emojiModule.toDBSaveEntities(o), $audioDetail.$src_content = $emojiModule.revertCommentTags($audioDetail.$src_content);
                    var i = $("#audioId").val();
                    $audioDetail.$isSending = !0;
                    var a = $("#title").val();
                    NetWork.request("insert_a_comment", {
                        type: 1,
                        recordid: i,
                        recordtitle: a,
                        content: o,
                        src_commentid: $audioDetail.$src_commentid,
                        src_only_for_group_post_userid: $audioDetail.$src_only_for_group_post_userid,
                        src_userid: $audioDetail.$src_userid,
                        src_content: $audioDetail.$src_content,
                        src_nickname: $audioDetail.$src_nickname
                    }, function (o) {
                        if ($audioDetail.$firtContent = !0, $audioDetail.$isSending = !1, null != o) {
                            $audioDetail.$comment_input_real.blur(), $audioDetail.$comment_bottom.hide(), $audioDetail.$comment_bottom_face.show(), null != $audioDetail.$droploadMe && $audioDetail.$droploadMe.clearAfterEmptyContenBox(), $audioDetail.$comment_input_real.val(""), $audioDetail.$comment_input_face.html("发评论"), $audioDetail.$comment_input_real.attr("placeholder", "发评论"), $audioDetail.clearReplyInfos(), $audioDetail.$commentsList.empty(), $audioDetail.$currentCommentPage = 1, $audioDetail.loadComment(null, !0);
                            var i = $(".title_hint_left"),
                                a = i.html();
                            a = a.replace("评论", ""), a = a.replace("(", ""), a = a.replace(")", ""), a = parseInt(a), a = a ? a + 1 : 1, i.html("评论(" + a + ")");
                            var t = $("#talk_count").text();
                            t = parseInt(t), t = t ? t + 1 : 1, $("#talk_count").text(t + "")
                        } else alert("评论失败,请重试"), $audioDetail.$isSending = !1, $audioDetail.$firtContent = !0
                    })
                }
            }
        }), $audioDetail.signLayer(), $audioDetail.message(), $audioDetail.replyClickArea()
    }, replyClick: function (o, i, a, t) {
        $audioDetail.$src_commentid = o, $audioDetail.$src_only_for_group_post_userid = "", $audioDetail.$src_userid = i, $audioDetail.$src_content = a, $audioDetail.$src_nickname = t;
        var e = "回复:" + $audioDetail.$src_nickname;
        $audioDetail.$comment_input_face.html(e), $audioDetail.$comment_input_real.attr("placeholder", $audioDetail.$comment_input_face.html())
    }, replyClickArea: function () {
        $("#comment_div").on("click", ".comment_item_content", function () {
            var o = $(this).data("commentid"),
                i = $(this).data("userid"),
                a = $(this).data("commnet"),
                t = decodeURI($(this).data("nickname"));
            /android/.test(navigator.userAgent.toLowerCase()) || (t = $emojiModule.toNormalEntities(t)), $audioDetail.replyClick(o, i, a, t)
        })
    }, clearReplyInfos: function () {
        $audioDetail.$src_commentid = "", $audioDetail.$src_only_for_group_post_userid = "", $audioDetail.$src_userid = "", $audioDetail.$src_content = "", $audioDetail.$src_nickname = ""
    }, loadMoreComment: function (o) {
        $audioDetail.loadComment(o, !1)
    }, loadComment: function (o, i) {
        var a = $("#audioId").val();
        $audioDetail.$isLoading || $audioDetail.$isSending || ($audioDetail.$isLoading = !0, NetWork.request("query_comments", {
            type: 1,
            page: $audioDetail.$currentCommentPage,
            record_id: a,
            load_by_hot: $audioDetail.$isOrderByHot ? "1" : "-1"
        }, function (a) {
            var t = null == o ? $audioDetail.$droploadMe : o;
            console.log(t);
            return null == a ? ($audioDetail.$isLoading = !1, t.resetload(), void $(".dropload-down").eq(-2).remove()) : (a = $emojiModule.toNormalEntities(a), $audioDetail.$isLoading = !1, null == a || "" == a || "nomoredata" == a ? null != t && (t.noData(), t.lock()) : ($audioDetail.$currentCommentPage++, null == o && $audioDetail.$commentsList.empty(), $audioDetail.$commentsList.append(a), i && window.scrollTo(0, $("#new_comments_title_hint").offset().top), a.split("comment_item_left_icon").length - 1 < 3 && null != t && (t.noData(), t.lock())), void(null != t && t.resetload()))
        }))
    }, initCommentList: function () {
        var o = $("#comment_div");
        if (0 != o.size()) {
            $(".content_wrapper").dropload({
                scrollArea: window,
                distance: 100,
                domDown: {
                    domClass: "dropload-down",
                    domRefresh: '<div class="dropload-refresh t2 c3 commentEnd">↑上拉加载更多</div>',
                    domLoad: '<div class="dropload-load t2 c3 commentEnd"><span class="loading"></span>加载中...</div>',
                    domNoData: '<div class="dropload-noData t2 c3 commentEnd">已加载完</div>'
                },
                autoLoad: !1,
                loadDownFn: function (o) {
                    $audioDetail.loadMoreComment(o)
                }, returnTheObjWhenBind: function (o) {
                    $audioDetail.$droploadMe = o
                }
            })
        }
    }, clickZan: function (o, i, a) {
        if ($("#availableInfo").val()) {
            var t = $(".zan_icon_" + o),
                e = $(".zan_num_" + o),
                n = "1";
            t.each(function (o, i) {
                i.src.indexOf("zan_blue") == -1 ? i.src = "/images/comment/zan_blue.png" : (i.src = "/images/comment/zan_white.png", n = "2")
            }), e.each(function (o, i) {
                var a = i.getAttribute("data-num"),
                    t = i.getAttribute("data-zanstate");
                "" == t || "0" == t || "2" == t ? (a = parseInt(a) + 1, i.setAttribute("data-num", a), i.setAttribute("data-zanstate", "1")) : "1" == t && (a = parseInt(a) - 1, i.setAttribute("data-num", a), i.setAttribute("data-zanstate", "2")), i.innerHTML = 0 == a ? "" : a
            });
            var u = $("#audioId").val();
            NetWork.request("comment_zan", {
                type: 1,
                record_id: u,
                comment_id: o,
                src_userid: i,
                src_comment_content: a,
                new_zan_state: n
            }, function (o) {})
        }
    }, checkForKeyBoard: function () {
        var o = window.innerHeight;
        if (o < $audioDetail.$maxHeight - 10) $audioDetail.$isKeyBoardShown = !0;
        else if (o >= $audioDetail.$maxHeight - 10 && $audioDetail.$isKeyBoardShown) {
            var i = $audioDetail.$comment_input_real.val();
            "" != i && null != i || (i = "发评论", $audioDetail.$comment_input_real.attr("placeholder", "发评论"), $audioDetail.clearReplyInfos()), $audioDetail.$comment_input_face.html(i), $("body").off("touchmove"), $audioDetail.$isKeyBoardShown = !1, $audioDetail.$scrollTop && !$("#comment_input_real").val() && ($("body").scrollTop($audioDetail.$scrollTop), $audioDetail.$scrollTop = 0)
        }
        o > $audioDetail.$maxHeight ? $audioDetail.$maxHeight = o : o < $audioDetail.$minHeight && ($audioDetail.$minHeight = o)
    }, signLayer: function () {
        var o = $("#sign_url_base64").val();
        o && ($audioDetail.$signButton = $("#sign_button"), $audioDetail.$signLayer = $("#sign_layer"), $audioDetail.$signLayerClose = $("#sign_layer_close"), $audioDetail.$signLayerBigImage = $("#sign_layer_big_image"), $audioDetail.$signButton.click(function () {
            var i = $("#userId").val();
            if (null == $audioDetail.$signImgUrl) {
                var a = window.location.hostname.indexOf(".h5.");
                a === -1 ? $audioDetail.$signImgUrl = "/" + window.APPID + "/get_daily_sign/" + o + "/" + i : $audioDetail.$signImgUrl = "/get_daily_sign/" + o + "/" + i;
                var t = location.href.split("?")[0],
                    e = t.split("/"),
                    n = e[e.length - 1];
                $audioDetail.$signImgUrl = $audioDetail.$signImgUrl + "/" + n, $audioDetail.$signLayerBigImage.attr("src", $audioDetail.$signImgUrl), $audioDetail.$signLayer.show()
            } else $audioDetail.$signLayer.show();
            $playStatus.clickDailySign()
        }), $audioDetail.$signLayerClose.click(function () {
            $audioDetail.$signLayer.hide()
        }))
    }, message: function () {
        $audioDetail.$messageImageShowing = !1, $("#audioPic_message_img").click(function () {
            if ($audioDetail.$messageImageShowing) return 0;
            $("#message_layer").show(), $(".audioPic_message_red").hide(), $audioDetail.$messageImageShowing = !0, scrollTo(0, 0), $("body").on("touchmove", function (o) {
                o.preventDefault()
            });
            var o = {
                showed: !0
            };
            NetWork.request("message_showed", o, function (o) {})
        }), $("#message_close").click(function () {
            $audioDetail.$messageImageShowing = !1, $("body").off("touchmove"), $("#message_layer").hide()
        })
    }, dealTopCommentsEmoji: function () {
        var o, i = $("#comments_list_tophot"),
            a = i.html();
        try {
            o = $emojiModule.toNormalEntities(a)
        } catch (t) {
            o = a
        }
        i.html(o), i.show()
    }, changeCommentOrderRole: function (o) {
        $audioDetail.$isOrderByHot && !o ? ($audioDetail.$isOrderByHot = !1, $("#new_order").addClass("c4"), $("#hot_order").removeClass("c4"), $audioDetail.$droploadMe.clearAfterEmptyContenBox(), $audioDetail.$currentCommentPage = 1, $audioDetail.loadMoreComment(null)) : !$audioDetail.$isOrderByHot && o && ($audioDetail.$isOrderByHot = !0, $("#hot_order").addClass("c4"), $("#new_order").removeClass("c4"), $audioDetail.$droploadMe.clearAfterEmptyContenBox(), $audioDetail.$currentCommentPage = 1, $audioDetail.loadMoreComment(null))
    }, deleteComment: function (o) {
        dialog({
            type: "confirm",
            desc: "评论删除后将无法在评论区显示,确认删除?",
            cancelTxt: "取消",
            confirmTxt: "确认",
            confirm: function () {
                NetWork.request("delete_comment", {
                    comment_id: o,
                    type: 1,
                    record_id: $audioFunction.$audioId
                }, function (i) {
                    if (console.log(i), i = JSON.parse(i), 0 == i.code) {
                        var a = "div[data-commentid='" + o + "']",
                            t = $(a);
                        t.remove();
                        var e = $(".title_hint_left"),
                            n = e.html();
                        n = n.replace("评论", ""), n = n.replace("(", ""), n = n.replace(")", ""), n = parseInt(n), n = n ? n - 1 : 0, e.html("评论(" + n + ")");
                        var u = $("#talk_count").text();
                        u = parseInt(u), u = u ? u - 1 : 0, $("#talk_count").text(u + "")
                    } else alert("删除失败")
                })
            }
        })
    }
}, $audioFunction = {
    status: {
        nextStatus: !0,
        prevStatus: !0,
        endStatus: !0,
        played: !1
    },
    $audio: null,
    $audioDom: $("#audioInfo"),
    $audioTryPlayTimes: null,
    $playButton: null,
    $orgAudioLength: 0,
    $timeLine: null,
    $timeBall: null,
    $progressWrapper: null,
    $panStartLeft: null,
    $startUpdate: null,
    $isEndEventSended: !1,
    $runTime: null,
    $waitForCheckAutoPlay: !0,
    $jiemaRetryFlag: !0,
    $net_error_try_times: 0,
    $RE_TRY_CDN: 2,
    $RE_TRY_OUR_SERVER: 1,
    $RE_NO_TRY_AGAIN: 0,
    $canConnectAudioSourseRetryFlag: -1,
    $audio_type: $("#audioInfo").data("type"),
    $audio_m3u8: $("#audioInfo").data("m3u8"),
    $audio_mp3: $("#audioInfo").data("mp3"),
    $orgin_mp3: $("#audioInfo").data("orginmp3"),
    $audio_quick: !1,
    $audio_quick_sign: !1,
    $audio_last_currentTime: 0,
    $audio_src_replace_times: 2,
    $audio_mode_timer: -1,
    $audioId: $("#audioId").val(),
    $productId: $("#productId").val(),
    init: function () {
        function o(o) {
            $(".custom_modal").html(o), $(".custom_modal").removeClass("custom_modal_hide"), $audioFunction.$audio_mode_timer = setTimeout(function () {
                $(".custom_modal").addClass("custom_modal_hide"), clearTimeout($audioFunction.$audio_mode_timer)
            }, 1500)
        }
        $audioFunction.$canConnectAudioSourseRetryFlag = $audioFunction.$RE_TRY_CDN, $audioFunction.$orgAudioLength = $("#mapPost").data("orgaudiolength"), $audioFunction.$audio = $("#audioInfo"), $audioFunction.$playButton = $("#playButton"), $audioFunction.$timeLine = $(".timeLine"), $audioFunction.$timeBall = $(".timeBall"), $audioFunction.$runTime = $(".runTime"), $audioFunction.$progressWrapper = $(".progressWrapper"), $audioFunction.$audioTime = $(".audioTime"), $audioFunction.initPlayMode(), document.getElementById("audioInfo").addEventListener("play", function () {
            "undefined" != typeof e_watcher && e_watcher.report_audio_success_play()
        }), document.getElementById("audioInfo").addEventListener("timeupdate", function () {
            "undefined" != typeof e_watcher && e_watcher.report_audio_timeupdate(1e3 * $audioFunction.$audio[0].currentTime)
        }), $audioFunction.$startUpdate = setInterval(function () {
            $audioFunction.updateStates()
        }, 500), $audioFunction.$monitorAbort = setInterval(function () {
            $audioFunction.$audio_last_currentTime - $audioFunction.$audio[0].currentTime > 1 || $audioFunction.$audio[0].currentTime - $audioFunction.$audio_last_currentTime > 1 ? $audioFunction.$audio_last_currentTime = $audioFunction.$audio[0].currentTime : $audioFunction.$audio_last_currentTime != $audioFunction.$audio[0].currentTime || 0 != $audioFunction.$audio[0].currentTime || $audioFunction.$audio[0].paused || 4 != $audioFunction.$audio[0].readyState || $audioFunction.$audio_src_replace_times > 0 && ($audioEvent.reloadAudioWhenError(!0), $audioFunction.$audio_src_replace_times--)
        }, 5e3), $audioFunction.$playButton.on("click", function () {
            $audioFunction.$waitForCheckAutoPlay = !1, "undefined" != typeof e_watcher && e_watcher.report_audio_manual_played(), $audioFunction.isAudioPlaying() ? ($audioFunction.$audio[0].pause(), $audioFunction.updateStates()) : ($audioFunction.$audio[0].play(), $audioFunction.updateStates())
        }), $audioFunction.panTimeBall(), $(".hrefClick").click(function () {
            $audioFunction.status.nextStatus && $(this).hasClass("next") && ($audioFunction.status.nextStatus = !1), $audioFunction.status.nextStatus && $(this).hasClass("prev") && ($audioFunction.status.prevStatus = !1);
            var o = $(this).attr("data-href");
            o && eJump(o, !0)
        }), $audioFunction.$progressWrapper.on("click", function (o) {
            if ($audioFunction.$audio[0].duration) {
                o.stopPropagation(), clearInterval($audioFunction.$startUpdate), clearInterval($audioFunction.$monitorAbort), $audioFunction.$audio[0].pause();
                var i = $audioFunction.$progressWrapper.width(),
                    a = 0;
                if (0 != i) {
                    var t = o.pageX - $audioFunction.$audioTime.width();
                    a = $audioFunction.$audio[0].duration * (t / i)
                }
                $audioFunction.$audio_quick = !0, $audioFunction.$audio_quick_sign = !0;
                try {
                    $audioFunction.$audio[0].currentTime = parseInt(a)
                } catch (e) {}
                $audioFunction.$audio[0].play(), $audioFunction.updateStates(), $audioFunction.$startUpdate = setInterval(function () {
                    $audioFunction.updateStates()
                }, 500), $audioFunction.$audio_last_currentTime = parseInt(a), $audioFunction.$monitorAbort = setInterval(function () {
                    $audioFunction.$audio_last_currentTime - $audioFunction.$audio[0].currentTime > 1 || $audioFunction.$audio[0].currentTime - $audioFunction.$audio_last_currentTime > 1 ? $audioFunction.$audio_last_currentTime = $audioFunction.$audio[0].currentTime : $audioFunction.$audio[0].paused || 4 != $audioFunction.$audio[0].readyState || (console.log("应该换资源了"), $audioFunction.$audio_src_replace_times > 0 && ($audioEvent.reloadAudioWhenError(!0), $audioFunction.$audio_src_replace_times--))
                }, 6e3)
            }
        }), $audioFunction.$audio.on("ended", function () {
            var o = $("#is_try").val(),
                i = $("#purchased").val(),
                a = $("#paymentUrl").val();
            if (1 == o && 0 == i && $audioFunction.$productId) return void dialog({
                type: "confirm",
                desc: "试听已结束，如需收看该专栏完整内容，请订阅专栏",
                cancelTxt: "取消",
                confirmTxt: "订阅专栏",
                confirm: function () {
                    setTimeout(function () {
                        eJump(a, !1)
                    }, 100)
                }
            });
            if ($("#audio_mode").hasClass("orderBackground")) {
                $audioFunction.status.nextStatus && ($audioFunction.status.nextStatus = !1);
                var t = $(".next").attr("data-href");
                t.indexOf("content") >= 0 ? setTimeout(function () {
                    eJump(t, !1)
                }, 1e3) : $audioFunction.$audio[0].pause()
            } else $("#audio_mode").hasClass("loopBackground") && ($audioFunction.$audio[0].play(), $audioFunction.$isEndEventSended || ($audioFunction.$isEndEventSended = !0), $audioFunction.updateStates());
            0 === $audioActionAnalyse.$type.length && ($("#audioType") && $("#audioType").val() ? $audioActionAnalyse.$type = 1 : $audioActionAnalyse.$type = 0), 0 === $audioActionAnalyse.$audioId.length && ($audioActionAnalyse.$audioId = $("#audioId").val())
        }), setTimeout(function () {
            $audioFunction.isAudioPlaying() || $audioFunction.playWhenInit()
        }, 1e3), $("#audioType") && $("#audioType").val() ? $audioActionAnalyse.$type = 1 : $audioActionAnalyse.$type = 0, $audioActionAnalyse.$audioId = $("#audioId").val(), $audioActionAnalyse.playCount(), $audioFunction.$audio.on("canplay", function () {
            $audioContinue.readTime($audioFunction.$audioId, $audioFunction.$audioDom[0]), $audioFunction.$audio.off("canplay")
        }), setTimeout(function () {
            $("#playButton").removeClass("loadingBackground")
        }, 8e3), $(".audioModeWrapper").click(function () {
            var i = $("#audio_mode"),
                a = "";
            $(i).hasClass("onlyOneBackground") ? ($(i).removeClass("onlyOneBackground"), $(i).addClass("orderBackground"), a = "orderBackground", o("列表循环")) : $(i).hasClass("orderBackground") ? ($(i).removeClass("orderBackground"), $(i).addClass("loopBackground"), a = "loopBackground", o("单曲循环")) : $(i).hasClass("loopBackground") && ($(i).removeClass("loopBackground"), $(i).addClass("onlyOneBackground"), a = "onlyOneBackground", o("单曲播放"));
            var t = $("#productId").val(),
                e = $("#audioId").val();
            t ? localStorage.setItem(t, a) : localStorage.setItem(e, a)
        })
    }, initPlayMode: function () {
        var o = localStorage.getItem($audioFunction.$productId),
            i = localStorage.getItem($audioFunction.$audioId);
        null != o ? ($("#audio_mode").get(0).className = "", $("#audio_mode").addClass(o)) : null != i && ($("#audio_mode").get(0).className = "", $("#audio_mode").addClass(i))
    }, playWhenInit: function () {
        if ($audioFunction.$audioTryPlayTimes++ > 3) return void($audioFunction.$waitForCheckAutoPlay = !1);
        if ($audioFunction.isAudioPlaying()) return void($audioFunction.$waitForCheckAutoPlay = !1);
        if ($audioFunction.$waitForCheckAutoPlay) {
            $audioFunction.$audio[0].pause(), $audioFunction.$audio[0].play();
            try {
                wx.getNetworkType({
                    success: function (o) {
                        $audioFunction.isAudioPlaying() || $audioFunction.$audio[0].play()
                    }
                })
            } catch (o) {}
            setTimeout(function () {
                $audioFunction.playWhenInit()
            }, 1e3), $audioFunction.loading(), setTimeout(function () {
                $("#playButton").removeClass("loadingBackground")
            }, 3e3)
        }
    }, updateStates: function () {
        $audioDetail.checkForKeyBoard();
        var o = $audioFunction.$audio[0].currentTime,
            i = $audioFunction.$audio[0].duration,
            a = $audioFunction.isAudioPlaying(),
            t = $audioFunction.$audio[0].readyState;
        if (0 == $audioFunction.$orgAudioLength && ($audioFunction.$orgAudioLength = i), $audioFunction.$audio_quick_sign) {
            if ($audioFunction.$audio_quick || 4 == t) {
                var e = 0 != i ? o / i : 0;
                e = 100 * e, $audioFunction.$timeLine.css("width", e + "%");
                $(".progressWrapper").width();
                e <= 100 ? $audioFunction.$timeBall.css("left", e + "%") : $audioFunction.$timeBall.css("left", "0"), $audioFunction.$runTime.html($audioFunction.timeChange(Math.floor(o))), $audioFunction.$audio_quick && ($audioFunction.$audio_quick = !1), 4 == t && ($audioFunction.$audio_quick_sign = !1)
            }
        } else {
            var e = 0 != i ? o / i : 0;
            e = 100 * e, $audioFunction.$timeLine.css("width", e + "%");
            $(".progressWrapper").width();
            e <= 100 ? $audioFunction.$timeBall.css("left", e + "%") : $audioFunction.$timeBall.css("left", "0"), $audioFunction.$runTime.html($audioFunction.timeChange(Math.floor(o)))
        }!a || 4 != t && 3 != t ? $audioFunction.$playButton.hasClass("pauseBackground") && $audioFunction.$playButton.removeClass("pauseBackground") : $audioFunction.$playButton.hasClass("pauseBackground") || $audioFunction.$playButton.addClass("pauseBackground"), $audioFunction.loading(), o > .01 && $audioContinue.readTime($audioFunction.$audioId, $audioFunction.$audioDom[0]), $audioContinue.saveTime($audioFunction.$audioId, $audioFunction.$audioDom[0]), i - o <= 30 && (0 === $audioActionAnalyse.$type.length && ($("#audioType") && $("#audioType").val() ? $audioActionAnalyse.$type = 1 : $audioActionAnalyse.$type = 0), 0 === $audioActionAnalyse.$audioId.length && ($audioActionAnalyse.$audioId = $("#audioId").val()), $audioActionAnalyse.finishCount())
    }, panTimeBall: function () {
        new Hammer($audioFunction.$timeBall[0], {
            domEvents: !0
        });
        var o, i, a = $audioFunction.$progressWrapper.width();
        $audioFunction.$timeBall.on("panstart", function (o) {
            clearInterval($audioFunction.$startUpdate), clearInterval($audioFunction.$monitorAbort), $audioFunction.$startUpdate = null;
            var i = $audioFunction.$timeBall.css("left");
            i.indexOf("px") != -1 ? $audioFunction.$panStartLeft = parseFloat(i) : i.indexOf("%") != -1 && ($audioFunction.$panStartLeft = parseFloat(i) / 100 * a)
        }), $audioFunction.$timeBall.on("pan", function (t) {
            i = t.gesture.deltaX, o = i + $audioFunction.$panStartLeft, o <= 0 ? o = 0 : o > a && (o = a);
            var e = $audioFunction.timeChange(Math.floor(o / a * $audioFunction.$audio[0].duration));
            $(".runTime").html(e), $audioFunction.$timeLine.css("width", o + "px"), $audioFunction.$timeBall.css("left", o + "px")
        }), $audioFunction.$timeBall.on("panend", function (i) {
            if ($audioFunction.$audio[0].duration) {
                $audioFunction.$audio[0].pause(), $audioFunction.$audio_quick = !0, $audioFunction.$audio_quick_sign = !0;
                var t = o / a * $audioFunction.$audio[0].duration;
                try {
                    $audioFunction.$audio[0].currentTime = t
                } catch (e) {}
                o / a < 1 && $audioFunction.$audio[0].play(), $audioFunction.updateStates(), $audioFunction.$startUpdate = setInterval(function () {
                    $audioFunction.updateStates()
                }, 500), $audioFunction.$audio_last_currentTime = o / a * $audioFunction.$audio[0].duration, $audioFunction.$monitorAbort = setInterval(function () {
                    $audioFunction.$audio_last_currentTime - $audioFunction.$audio[0].currentTime > 1 || $audioFunction.$audio[0].currentTime - $audioFunction.$audio_last_currentTime > 1 ? $audioFunction.$audio_last_currentTime = $audioFunction.$audio[0].currentTime : $audioFunction.$audio[0].paused || 4 != $audioFunction.$audio[0].readyState || (console.log("应该换资源了"), $audioFunction.$audio_src_replace_times > 0 && ($audioEvent.reloadAudioWhenError(!0), $audioFunction.$audio_src_replace_times--))
                }, 6e3)
            }
        })
    }, isAudioPlaying: function () {
        var o = !1;
        return !$audioFunction.$audio[0].paused && $audioFunction.$audio[0].currentTime > 0 && (o = !0), o
    }, timeChange: function (o) {
        var i = parseInt(o / 60);
        i < 10 && (i = "0" + i);
        var a = parseInt(o % 60);
        return a < 10 && (a = "0" + a), isNaN(i) && (i = "00"), isNaN(a) && (a = "00"), i + ":" + a
    }, postAudioPlayError: function (o, i) {
        $.post("/report_audio_problem", {
            audio_id: $("#audioId").val(),
            audio_state_code: o,
            error_audio_url: $audioFunction.$audio[0].src,
            extra: i
        }, function () {})
    }, loading: function () {
        var o = $audioFunction.$audioDom[0].readyState;
        $audioFunction.$audio[0].ended ? $("#playButton").removeClass("loadingBackground") : 4 === o ? setTimeout(function () {
            $("#playButton").removeClass("loadingBackground")
        }, 200) : $audioFunction.$audio[0].paused || 3 != o ? 0 !== o && 1 !== o && 2 !== o && 3 !== o || $("#playButton").addClass("loadingBackground") : $("#playButton").removeClass("loadingBackground")
    }, moveFaskClick: function () {
        navigator.userAgent.toLowerCase().indexOf("iphone") >= 0 && 320 === window.innerWidth && $("#playButton").addClass("needsclick")
    }
}, $playStatus = {
    play: function () {}, pause: function () {}, prev: function () {}, next: function () {}, playNum: function () {}, end: function () {}, timelind: function () {}, dragtimeline: function () {}, share: function () {
        $playStatus.postToServer("share_count")
    }, clickDailySign: function () {
        $playStatus.postToServer("click_sign_count")
    }, postToServer: function (o) {
        var i = {
            audio_id: $("#audioId").val()
        };
        NetWork.request("audio_analyse/" + o, i, function (o) {})
    }
}, $audioEvent = {
    $net_error_try_times: 0,
    $decode_error_try_times: 0,
    $url_error_try_times: 0,
    $jiemaRetryFlag: !0,
    init: function () {
        $audioFunction.$audioDom[0].addEventListener("error", function () {
            var o = this.error.code,
                i = 1,
                a = 2,
                t = 3,
                e = 4;
            console.log("报错了,错误码是:" + o);
            var n = {
                limitTag: !1,
                errorEvent: "error",
                errorCode: o,
                srcType: 2,
                srcHref: $audioFunction.$audioDom.attr("src")
            };
            switch ($errorReprot.report(n), o) {
            case i:
                break;
            case a:
                $audioEvent.$net_error_try_times++, setTimeout(function () {
                    $audioEvent.reloadAudioWhenError(!1)
                }, 3e3);
                break;
            case t:
                $audioEvent.$decode_error_try_times++, $audioEvent.$decode_error_try_times <= 3 && ($audioEvent.$jiemaRetryFlag ? ($audioEvent.$jiemaRetryFlag = !1, $audioEvent.reloadAudioWhenError(!1)) : $audioEvent.reloadAudioWhenError(!0));
                break;
            case e:
                $audioEvent.$url_error_try_times <= 7 && ($audioEvent.$url_error_try_times++, $audioEvent.reloadAudioWhenError(!0))
            }
        })
    }, reloadAudioWhenError: function (o) {
        var i = $audioFunction.$audio[0].currentTime;
        "undefined" != typeof e_watcher && e_watcher.report_audio_source_changed(), o ? 1 == $audioFunction.$audio_type ? ($audioFunction.$audio[0].pause(), $audioFunction.$audio[0].src = $audioFunction.$audio_mp3, $audioFunction.$audio[0].load(), $audioFunction.$audio[0].play(i), $audioFunction.$audio_type = 0) : 0 == $audioFunction.$audio_type ? ($audioFunction.$audio[0].pause(), $audioFunction.$audio[0].src = $audioFunction.$orgin_mp3, $audioFunction.$audio[0].load(), $audioFunction.$audio[0].play(i), $audioFunction.$audio_type = -1) : $audioFunction.$audio_type == -1 ? $.ajax({
            type: "POST",
            url: $("#audioResourceServer").val(),
            data: {
                bizData: {
                    audioMp3Url: $audioFunction.$audio_mp3
                }
            },
            success: function (o) {
                $audioFunction.$audio[0].pause(), $audioFunction.$audio[0].src = o, $audioFunction.$audio[0].load(), $audioFunction.$audio[0].play(i), $audioFunction.$audio_type = -2
            }, error: function () {}
        }) : $audioFunction.$audio_type == -2 && ($audioFunction.$audio[0].pause(), $audioFunction.$audio[0].load(), $audioFunction.$audio[0].play(i)) : ($audioFunction.$audio[0].pause(), $audioFunction.$audio[0].load(), $audioFunction.$audio[0].play(i))
    }, onemptied: function () {
        NetWork.request("/report_audio_error/onemptied", {}, function (o) {}), $audioEvent.tryContinuePlay()
    }, onstalled: function () {
        NetWork.request("/report_audio_error/onstalled", {}, function (o) {}), $audioEvent.tryContinuePlay()
    }, isTrying: !1,
    tryContinuePlay: function () {
        $audioEvent.isTrying || ($audioEvent.isTrying = !0, setTimeout(function () {
            $audioFunction.isAudioPlaying() || $audioEvent.reloadAudioWhenError(), $audioEvent.isTrying = !1
        }, 2e3))
    }
};
var $ImageTextDesc = {
    init: function () {
        0 == $(".desc-mb").size() && ($("#detail_div img").each(function () {
            if ($(this).attr("alt")) {
                for (var o, i = $("#detail_div img"), a = [], t = 0; o = i[t++];) $(o).attr("alt") && (a = a.concat($(o).attr("src")));
                $(this).click(function (o) {
                    var i = $(this).attr("src");
                    wx.previewImage({
                        current: i,
                        urls: a
                    })
                })
            }
        }), $("#detail_div a img").off("click"))
    }
};
$(document).ready(function () {
    $hideDetail.init(), ($("#availableInfo").val() || $("#canTry").val()) && ($audioFunction.init(), $audioEvent.init()), $("#availableInfo").val() ? $audioDetail.init() : $payButton.init(), localStorage.setItem("needInitPage", !0), $emojiModule.init("comment_input_real"), $("#giftBuy").val() && $giftBuy.init();
    var o = $("#h5LogId").val();
    "undefined" != typeof e_watcher && e_watcher.e_init(o, e_watcher.TYPE_SOURCE_AUDIO), $subscribeModule.init(), $ImageTextDesc.init()
});