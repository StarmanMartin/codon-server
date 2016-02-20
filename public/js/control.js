$(function () {
    function clickHandler($this) {
        $this.closest('.nav-custem').find('li').removeClass('active');
        $this.parent().addClass('active');
        var $swContainer = $this.closest('.sw-container');
        $swContainer.find('.sw-box').removeClass('active');
        $swContainer.find('.sw-' + $this.data('linktype')).addClass('active');
    }

    function initContainer() {
        $('.sw-container .nav-custem a').click(function(){
            var $this = $(this);
            clickHandler($this);
        });
        
        $('.sw-container .nav-custem li.active').each(function(){
            var $this = $(this).find('a');
            clickHandler($this);
        });
    }

    initContainer();
});