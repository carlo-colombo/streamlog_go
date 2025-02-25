import {Component, OnInit} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {SseClient} from 'ngx-sse-client';
import {HttpHeaders} from '@angular/common/http';
import {NgFor} from '@angular/common';
import {DomSanitizer, SafeHtml} from '@angular/platform-browser';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, NgFor],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'app';
  lines: SafeHtml[] = []

  constructor(private sseClient: SseClient, private sanitizer: DomSanitizer) {
    const headers = new HttpHeaders().set('Authorization', `Basic YWRtaW46YWRtaW4=`);

    this.sseClient.stream('/logs?sse', {
      keepAlive: true,
      reconnectionDelay: 1_000,
      responseType: 'event'
    }, {headers}, 'GET')
      .subscribe((event) => {
        if (event.type === 'error') {
          const errorEvent = event as ErrorEvent;
          console.error(errorEvent.error, errorEvent.message);
        } else {
          const messageEvent = event as MessageEvent;

          this.lines.unshift(this.sanitizer.bypassSecurityTrustHtml( messageEvent.data));

          console.info(`SSE request with type "${messageEvent.type}" and data "${messageEvent.data}"`);
        }
      });
  }

  ngOnInit(): void {
    // throw new Error('Method not implemented.');
  }
}
