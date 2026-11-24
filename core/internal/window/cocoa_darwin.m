#import <Cocoa/Cocoa.h>
#import <QuartzCore/QuartzCore.h>

#include "cocoa_darwin.h"
#include "_cgo_export.h"

@interface PetPanel : NSPanel
@end

@implementation PetPanel
- (BOOL)canBecomeKeyWindow {
	return NO;
}
- (BOOL)canBecomeMainWindow {
	return NO;
}
@end

@interface PetView : NSView {
	NSPoint grab;
	NSPoint origin;
	BOOL moved;
}
@end

@implementation PetView
- (BOOL)acceptsFirstMouse:(NSEvent *)event {
	return YES;
}
- (BOOL)mouseDownCanMoveWindow {
	return NO;
}
- (void)mouseDown:(NSEvent *)event {
	grab = [NSEvent mouseLocation];
	origin = self.window.frame.origin;
	moved = NO;
}
- (void)mouseDragged:(NSEvent *)event {
	NSPoint p = [NSEvent mouseLocation];
	if (p.x == grab.x && p.y == grab.y) {
		return;
	}
	[self.window setFrameOrigin:NSMakePoint(origin.x + p.x - grab.x, origin.y + p.y - grab.y)];
	moved = YES;
}
- (void)mouseUp:(NSEvent *)event {
	if (!moved) {
		return;
	}
	moved = NO;
	NSUInteger idx = [[NSScreen screens] indexOfObject:self.window.screen];
	NSPoint o = self.window.frame.origin;
	goPetMoved(o.x, o.y, idx == NSNotFound ? 0 : (int)idx);
}
@end

static PetPanel *panel;
static NSMutableArray *images;

void petInit(void) {
	@autoreleasepool {
		[NSApplication sharedApplication];
		[NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
	}
}

int petScreen(int i, double *x, double *y, double *w, double *h, double *scale) {
	@autoreleasepool {
		NSArray *screens = [NSScreen screens];
		int n = (int)screens.count;
		if (i < 0 || i >= n) {
			return n;
		}
		NSScreen *s = screens[i];
		NSRect f = s.visibleFrame;
		*x = f.origin.x;
		*y = f.origin.y;
		*w = f.size.width;
		*h = f.size.height;
		*scale = s.backingScaleFactor;
		return n;
	}
}

bool petCreate(double x, double y, double w, double h, double scale, bool topmost,
               const uint8_t *pix, int frames, int pw, int ph) {
	@autoreleasepool {
		images = [[NSMutableArray alloc] initWithCapacity:frames];
		CGColorSpaceRef cs = CGColorSpaceCreateWithName(kCGColorSpaceSRGB);
		size_t n = (size_t)pw * ph * 4;
		for (int i = 0; i < frames; i++) {
			CFDataRef data = CFDataCreate(NULL, pix + i * n, n);
			CGDataProviderRef dp = CGDataProviderCreateWithCFData(data);
			CGImageRef img = CGImageCreate(pw, ph, 8, 32, pw * 4, cs,
			                               kCGImageAlphaPremultipliedLast | kCGBitmapByteOrderDefault,
			                               dp, NULL, false, kCGRenderingIntentDefault);
			CGDataProviderRelease(dp);
			CFRelease(data);
			if (img == NULL) {
				CGColorSpaceRelease(cs);
				return false;
			}
			[images addObject:(id)img];
			CGImageRelease(img);
		}
		CGColorSpaceRelease(cs);

		panel = [[PetPanel alloc] initWithContentRect:NSMakeRect(x, y, w, h)
		                                    styleMask:NSWindowStyleMaskBorderless | NSWindowStyleMaskNonactivatingPanel
		                                      backing:NSBackingStoreBuffered
		                                        defer:NO];
		panel.opaque = NO;
		panel.backgroundColor = [NSColor clearColor];
		panel.hasShadow = NO;
		panel.hidesOnDeactivate = NO;
		panel.releasedWhenClosed = NO;
		panel.level = topmost ? NSFloatingWindowLevel : NSNormalWindowLevel;
		panel.collectionBehavior = NSWindowCollectionBehaviorCanJoinAllSpaces |
		                           NSWindowCollectionBehaviorStationary |
		                           NSWindowCollectionBehaviorIgnoresCycle |
		                           NSWindowCollectionBehaviorFullScreenAuxiliary;
		panel.title = @"Claude Pet";

		PetView *view = [[PetView alloc] initWithFrame:NSMakeRect(0, 0, w, h)];
		view.wantsLayer = YES;
		view.layer.contentsGravity = kCAGravityResize;
		view.layer.contentsScale = scale;
		panel.contentView = view;
		[view release];
		petShowFrame(0);
		[panel orderFrontRegardless];
		return true;
	}
}

void petShowFrame(int i) {
	if (i < 0 || i >= (int)images.count) {
		return;
	}
	[CATransaction begin];
	[CATransaction setDisableActions:YES];
	panel.contentView.layer.contents = images[i];
	[CATransaction commit];
}

void petRun(double interval) {
	@autoreleasepool {
		NSTimer *t = [NSTimer timerWithTimeInterval:interval
		                                    repeats:YES
		                                      block:^(NSTimer *timer) {
			                                    goPetTick();
		                                      }];
		[[NSRunLoop mainRunLoop] addTimer:t forMode:NSRunLoopCommonModes];
		[NSApp run];
	}
}
