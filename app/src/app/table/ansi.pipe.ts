import { Pipe, PipeTransform } from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { FancyAnsi } from 'fancy-ansi';

@Pipe({
  name: 'ansi',
  standalone: true
})
export class AnsiPipe implements PipeTransform {
  private ansi = new FancyAnsi();

  constructor(private sanitizer: DomSanitizer) {}

  transform(value: string): SafeHtml {
    if (!value) return '';
    const html = this.ansi.toHtml(value);
    return this.sanitizer.bypassSecurityTrustHtml(html);
  }
} 