const picture = document.getElementsByClassName('zoom');
const ele = document.getElementsByClassName('ele');
const prev = document.getElementsByClassName('prev');
const next = document.getElementsByClassName('next');
const color = document.getElementsByClassName('pc');
const name_color = document.getElementById('name-color');

for (let i = 0; i < color.length; i++) {
    
    color[i].addEventListener('click', () =>{
        name_color.innerHTML = color[i].classList[1]
        for(var j = 0; j < color.length; j++) {
            color[j].innerHTML = ""
        }
        color[i].innerHTML = '<i class="fa-solid fa-check fa-xs" style="color: #ffffff;"></i>'
    });
}

max = ele.length-1
number = 0
for (let i = 0; i < ele.length; i++) {
    ele[i].addEventListener('click', () => {
        // console.log(ele[i].innerHTML);
        ele[i].classList.add('clicked')
        ele[number].classList.remove('clicked')
        number = i
        picture[0].innerHTML = ele[i].innerHTML
    });
}

prev[0].addEventListener('click', () => {
    console.log('Previous')
    if (max<2) {return}
    if (number > 0) {
        
        number = number - 1
        ele[number].classList.add('clicked')
        ele[number+1].classList.remove('clicked')
        picture[0].innerHTML = ele[number].innerHTML
    }else {
        
        number = max
        ele[number].classList.add('clicked')
        ele[0].classList.remove('clicked')
        picture[0].innerHTML = ele[number].innerHTML
    }
    
});

next[0].addEventListener('click', () => {
    console.log('next click')
    if (max<2) {return}
    if (number < max) {
        
        number = number + 1
        ele[number].classList.add('clicked')
        ele[number-1].classList.remove('clicked')
        picture[0].innerHTML = ele[number].innerHTML
    }else {
        
        number = 0
        ele[number].classList.add('clicked')
        ele[max].classList.remove('clicked')
        picture[0].innerHTML = ele[number].innerHTML
    }
});
