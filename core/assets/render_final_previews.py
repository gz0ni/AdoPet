from pathlib import Path
from PIL import Image, ImageDraw
import json

ROOT=Path(__file__).resolve().parent
atlas=Image.open(ROOT/'final/spritesheet-extended.png').convert('RGBA')
states=['idle','running-right','running-left','waving','jumping','failed','waiting','running','review']
durations=[[280,110,110,140,140,320],[120]*7+[220],[120]*7+[220],[140]*3+[280],[140]*4+[280],[140]*7+[240],[150]*5+[260],[120]*5+[220],[150]*5+[280]]
out=ROOT/'qa/final-previews';out.mkdir(parents=True,exist_ok=True)
frames_root=out/'frames';frames_root.mkdir(exist_ok=True)
all_frames=[];all_durations=[]
def cell(r,c): return atlas.crop((192*c,208*r,192*(c+1),208*(r+1)))
def board(im,label):
    b=Image.new('RGB',(384,456),'#c8cbd2')
    f=im.resize((384,416),Image.Resampling.NEAREST)
    b.paste(f,(0,40),f);ImageDraw.Draw(b).text((12,12),label,fill='#16171e')
    return b
def gif(frames,ds,path):
    frames[0].save(path,save_all=True,append_images=frames[1:],duration=ds,loop=0,disposal=2,optimize=False)
for r,(state,ds) in enumerate(zip(states,durations)):
    target=frames_root/state;target.mkdir(exist_ok=True)
    fs=[cell(r,c) for c in range(len(ds))]
    for c,f in enumerate(fs): f.save(target/f'{c:02}.png')
    all_frames.extend(board(f,state) for f in fs);all_durations.extend(ds)
gif(all_frames,all_durations,out/'all-states.gif')
transition=[cell(0,0)]+[cell(4,c) for c in range(5)]+[cell(0,0)]
gif(transition,[500,140,140,140,140,280,500],out/'idle-jump-idle.gif')
directions=['000','022.5','045','067.5','090','112.5','135','157.5','180','202.5','225','247.5','270','292.5','315','337.5']
look=[board(cell(9+i//8,i%8),'look '+d) for i,d in enumerate(directions)]
gif(look,[180]*16,out/'look-loop.gif')
stills=Image.new('RGB',(384*2,456*2),'#c8cbd2')
for i,(r,c) in enumerate([(0,0),(3,2),(4,2),(7,2)]): stills.paste(board(cell(r,c),states[r]),((i%2)*384,(i//2)*456))
stills.save(out/'four-state-stills.png')
video=out/'video-frames';video.mkdir(exist_ok=True)
n=0
for f,ms in zip(all_frames,all_durations):
    for _ in range(max(1,round(ms/100))):
        f.save(video/f'{n:04}.png');n+=1
print(json.dumps({'ok':True,'states':9,'video_frames':n,'source':str(ROOT/'final/spritesheet-extended.png')}))
